package handler

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/useCase"
)

func RegisterHTTPEndPoints(u useCase.ITask, logger *slog.Logger, router *chi.Mux) *ITask {
	handler := NewHandler(u, logger)
	router.Route("/v1/task", func(router chi.Router) {
		router.Get("/{taskID}", handler.FindOne)
		router.Post("/", handler.Create)
		router.Put("/", handler.Update)
		router.Delete("/{taskID}", handler.Delete)
	})
	return handler
}
