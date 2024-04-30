package server

import (
	"context"
	"log"

	taskHandler "github.com/henriqueassiss/advanced-golang-api/internal/domain/task/handler"
	taskMock "github.com/henriqueassiss/advanced-golang-api/internal/domain/task/mock"
	taskRepository "github.com/henriqueassiss/advanced-golang-api/internal/domain/task/repository"
	taskUseCase "github.com/henriqueassiss/advanced-golang-api/internal/domain/task/useCase"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/errorMsg"
	"github.com/jwalton/gchalk"
)

func (s *Server) InitDomains() {
	log.Println(gchalk.Yellow("Domain: starting"))

	s.initAuthentication()

	log.Println(gchalk.Blue("Domain: done"))

	if s.cfg.App.Environment == "development" {
		s.InitMocks()
	}
}

func (s *Server) InitMocks() {
	log.Println(gchalk.Yellow("Mock: starting"))

	err := taskMock.Mock(s.sqlx)
	if err != nil && err != errorMsg.ErrTableIsPopulated {
		GracefulShutdown(context.Background(), s)
	}

	log.Println(gchalk.Blue("Mock: done"))
}

func (s *Server) initAuthentication() {
	newTaskRepo := taskRepository.New(s.sqlx)
	newTaskUseCase := taskUseCase.New(newTaskRepo, s.logger, s.cache)
	taskHandler.RegisterHTTPEndPoints(newTaskUseCase, s.logger, s.router)
}
