package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/henriqueassiss/advanced-golang-api/config"
	"github.com/henriqueassiss/advanced-golang-api/internal/middleware"
	"github.com/henriqueassiss/advanced-golang-api/third_party/cache"
	"github.com/henriqueassiss/advanced-golang-api/third_party/logger"

	db "github.com/henriqueassiss/advanced-golang-api/third_party/database"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/jwalton/gchalk"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
)

type Server struct {
	Version string
	cfg     *config.Config
	logger  *slog.Logger

	cache *redis.Client
	sqlx  *sqlx.DB

	cors   *cors.Cors
	router *chi.Mux

	httpServer *http.Server
}

type Options func(opts *Server) error

func New(opts ...Options) *Server {
	s := defaultServer()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return s
}

func WithVersion(version string) Options {
	return func(opts *Server) error {
		log.Println(gchalk.Green("Starting API version:", version))
		opts.Version = version
		return nil
	}
}

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(false, false),
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
	s.newLogger()
	s.setCors()
	s.newCache()
	s.newDatabase()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitDomains()
}

func (s *Server) newLogger() {
	log.Println(gchalk.Yellow("Logger: starting"))

	s.logger = logger.New()

	log.Println(gchalk.Blue("Logger: done"))
}

func (s *Server) setCors() {
	log.Println(gchalk.Yellow("Cors: setting"))

	s.cors = cors.New(
		cors.Options{
			AllowedOrigins: s.cfg.Cors.AllowedOrigins,
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})

	log.Println(gchalk.Blue("Cors: setted"))
}

func (s *Server) newCache() {
	log.Println(gchalk.Yellow("Cache: starting"))

	if s.cfg.Cache.Address == "" {
		s.logger.Error("Please, fill in redis credentials in .env file or set in environment variable")
		GracefulShutdown(context.Background(), s)
	}

	s.cache = cache.New(s.cfg.Cache)

	log.Println(gchalk.Blue("Cache: done"))
}

func (s *Server) newDatabase() {
	log.Println(gchalk.Yellow("Database: starting"))

	if s.cfg.Database.Driver == "" {
		log.Println(gchalk.Red("Please, fill in database credentials in .env file or set in environment variable"))
		GracefulShutdown(context.Background(), s)
	}

	db, err := db.NewSqlx(s.cfg.Database)
	if err != nil {
		s.logger.Error(err.Error())
		GracefulShutdown(context.Background(), s)
	}

	s.sqlx = db
	s.sqlx.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.sqlx.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.sqlx.SetConnMaxLifetime(s.cfg.Database.ConnectionsMaxLifeTime)

	log.Println(gchalk.Blue("Database: done"))
}

func (s *Server) newRouter() {
	log.Println(gchalk.Yellow("Router: starting"))

	s.router = chi.NewRouter()

	log.Println(gchalk.Blue("Router: done"))
}

func (s *Server) setGlobalMiddleware() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})
	s.router.Use(s.cors.Handler)
	s.router.Use(middleware.Json)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
}

func (s *Server) Run() {
	log.Println(gchalk.Yellow("Server: starting"))

	s.httpServer = &http.Server{
		Addr:              s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Api.ReadHeaderTimeout,
	}

	go func() {
		start(s)
	}()

	log.Println(gchalk.Blue("Server: done"))

	_ = GracefulShutdown(context.Background(), s)
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func start(s *Server) {
	log.Println(gchalk.Green(fmt.Sprintf("Serving at %s:%s", s.cfg.Api.Host, s.cfg.Api.Port)))

	var err error
	if s.cfg.App.Environment == "development" {
		err = s.httpServer.ListenAndServe()
	} else {
		err = s.httpServer.ListenAndServeTLS(s.cfg.App.CertDir, s.cfg.App.KeyDir)
	}

	if err != nil {
		s.logger.Error(err.Error())
		GracefulShutdown(context.Background(), s)
	}
}

func GracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println(gchalk.Red("Server: shutting down"))

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Error(err.Error())
	}

	return nil
}
