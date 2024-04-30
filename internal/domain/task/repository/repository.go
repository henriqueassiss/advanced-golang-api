package repository

import (
	"context"
	"strings"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/schema"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ITask interface {
	FindOne(ctx context.Context, params schema.QueryParams) (*task.Schema, error)
	FindMany(ctx context.Context, params schema.QueryParams) ([]task.Schema, error)
	Create(ctx context.Context, t *task.Schema) error
	Update(ctx context.Context, t *task.Schema) error
	Delete(ctx context.Context, taskID uint64) error
}

type Task struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Task {
	return &Task{
		db: db,
	}
}

func (r *Task) FindOne(ctx context.Context, params schema.QueryParams) (*task.Schema, error) {
	if params.Select == "" {
		params.Select = "t.*"
	}

	query := schema.PrepareFindQuery(Select, params)

	var t task.Schema
	err := r.db.GetContext(ctx, &t, query)

	return &t, err
}

func (r *Task) FindMany(ctx context.Context, params schema.QueryParams) ([]task.Schema, error) {
	if params.Select == "" {
		params.Select = "t.*"
	}

	query := schema.PrepareFindQuery(Select, params)

	var ts []task.Schema
	err := r.db.SelectContext(ctx, &ts, query)

	return ts, err
}

func (r *Task) Create(ctx context.Context, t *task.Schema) error {
	fields, values := schema.ParseFieldsToInsertQuery(t)

	query := strings.Replace(InsertInto, "?", fields, 1)

	query = strings.Replace(query, "?", values, 1)

	_, err := r.db.ExecContext(ctx, query)

	return err
}

func (r *Task) Update(ctx context.Context, t *task.Schema) error {
	fields := schema.ParseFieldsToUpdateQuery(t, "id", "task_colors", "task_infos")

	query := strings.Replace(Update, "?", fields, 1)

	_, err := r.db.ExecContext(ctx, query, t.ID)

	return err
}

func (r *Task) Delete(ctx context.Context, taskID uint64) error {
	_, err := r.db.ExecContext(ctx, Delete, taskID)

	return err
}
