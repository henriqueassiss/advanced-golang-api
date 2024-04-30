package useCase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/repository"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/schema"
	"github.com/redis/go-redis/v9"
)

type ITask interface {
	FindOne(ctx context.Context, taskID uint64) (*task.Schema, error)
	Create(ctx context.Context, t *task.Schema) error
	Update(ctx context.Context, t *task.Schema) error
	Delete(ctx context.Context, taskID uint64) error
}

type Task struct {
	repository repository.ITask
	logger     *slog.Logger
	cache      *redis.Client
}

func New(repo repository.ITask, logger *slog.Logger, cache *redis.Client) *Task {
	return &Task{
		repository: repo,
		logger:     logger,
		cache:      cache,
	}
}

func (uc *Task) FindOne(ctx context.Context, taskID uint64) (*task.Schema, error) {
	return uc.repository.FindOne(ctx, schema.QueryParams{
		Where: fmt.Sprintf("t.id = %d", taskID),
	})
}

func (uc *Task) Create(ctx context.Context, t *task.Schema) error {
	return uc.repository.Create(ctx, t)
}

func (uc *Task) Update(ctx context.Context, t *task.Schema) error {
	return uc.repository.Update(ctx, t)
}

func (uc *Task) Delete(ctx context.Context, taskID uint64) error {
	_, err := uc.repository.FindOne(ctx, schema.QueryParams{
		Select: "c.id",
		Where:  fmt.Sprintf("c.id = %d", taskID),
	})
	if err != nil {
		return err
	}

	err = uc.repository.Delete(ctx, taskID)

	return err
}
