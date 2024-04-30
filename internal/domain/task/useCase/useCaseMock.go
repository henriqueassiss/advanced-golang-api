package useCase

import (
	"context"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
)

type TaskMock struct {
	FindOneFunc func(ctx context.Context, taskID uint64) (*task.Schema, error)
	CreateFunc  func(ctx context.Context, t *task.Schema) error
	UpdateFunc  func(ctx context.Context, t *task.Schema) error
	DeleteFunc  func(ctx context.Context, taskID uint64) error
}

func (uc *TaskMock) FindOne(ctx context.Context, taskID uint64) (*task.Schema, error) {
	return uc.FindOneFunc(ctx, taskID)
}

func (uc *TaskMock) Create(ctx context.Context, t *task.Schema) error {
	return uc.CreateFunc(ctx, t)
}

func (uc *TaskMock) Update(ctx context.Context, t *task.Schema) error {
	return uc.UpdateFunc(ctx, t)
}

func (uc *TaskMock) Delete(ctx context.Context, taskID uint64) error {
	return uc.DeleteFunc(ctx, taskID)
}
