package useCase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/repository"

	"github.com/henriqueassiss/advanced-golang-api/third_party/cache"
	"github.com/henriqueassiss/advanced-golang-api/third_party/database"
	"github.com/henriqueassiss/advanced-golang-api/third_party/logger"

	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTaskUseCase_FindOne(t *testing.T) {
	logger := logger.New()
	db, mock := database.NewSqlxMock(t)
	cacheMock := cache.NewMock(t)
	r := repository.New(db)
	uc := New(r, logger, cacheMock)
	defer db.Close()

	type args struct {
		ctx    context.Context
		taskID uint64
	}

	type want struct {
		t   *task.Schema
		err error
	}

	type test struct {
		name string
		args
		beforeTest func(uint64)
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				ctx:    context.TODO(),
				taskID: 1,
			},
			beforeTest: func(taskID uint64) {
				taskRows := mock.NewRows([]string{
					"id",
					"title",
					"description",
					"updated_at",
				}).AddRow(
					taskID,
					"Test",
					"Test",
					sql.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
					},
				)
				mock.ExpectQuery("SELECT t.* FROM tasks t WHERE t.id =").WillReturnRows(taskRows)
			},
			want: want{
				t: &task.Schema{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					UpdatedAt: sql.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(tt.args.taskID)

			got, err := uc.FindOne(tt.args.ctx, tt.args.taskID)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.t, got)
		})
	}
}

func TestTaskUseCase_Create(t *testing.T) {
	logger := logger.New()
	db, mock := database.NewSqlxMock(t)
	cacheMock := cache.NewMock(t)
	r := repository.New(db)
	uc := New(r, logger, cacheMock)
	defer db.Close()

	type args struct {
		ctx context.Context
		t   *task.Schema
	}

	type want struct {
		err error
	}

	type test struct {
		name string
		args
		beforeTest func(*task.Schema)
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				ctx: context.TODO(),
				t: &task.Schema{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					UpdatedAt: sql.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
					},
				},
			},
			beforeTest: func(t *task.Schema) {
				mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlxmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(tt.args.t)

			err := uc.Create(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestTaskUseCase_Update(t *testing.T) {
	logger := logger.New()
	db, mock := database.NewSqlxMock(t)
	cacheMock := cache.NewMock(t)
	r := repository.New(db)
	uc := New(r, logger, cacheMock)
	defer db.Close()

	type args struct {
		ctx context.Context
		t   *task.Schema
	}

	type want struct {
		err error
	}

	type test struct {
		name string
		args
		beforeTest func(*task.Schema)
		want
	}

	tests := []test{
		{
			name: "Success - No images",
			args: args{
				ctx: context.TODO(),
				t: &task.Schema{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					UpdatedAt: sql.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
					},
				},
			},
			beforeTest: func(t *task.Schema) {
				mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlxmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(tt.args.t)

			err := uc.Update(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestTaskUseCase_Delete(t *testing.T) {
	logger := logger.New()
	db, mock := database.NewSqlxMock(t)
	cacheMock := cache.NewMock(t)
	r := repository.New(db)
	uc := New(r, logger, cacheMock)
	defer db.Close()

	type args struct {
		ctx    context.Context
		taskID uint64
	}

	type want struct {
		err error
	}

	type test struct {
		name string
		args
		beforeTest func(uint64)
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				ctx:    context.TODO(),
				taskID: 1,
			},
			beforeTest: func(taskID uint64) {
				taskRows := mock.NewRows([]string{"id"}).AddRow(taskID)
				mock.ExpectQuery("SELECT (.+) FROM tasks").WillReturnRows(taskRows)

				mock.ExpectExec("DELETE FROM tasks").WithArgs(taskID).WillReturnResult(sqlxmock.NewResult(0, 1))
			},
		},
		{
			name: "Fail - Invalid task",
			args: args{
				ctx:    context.TODO(),
				taskID: 0,
			},
			beforeTest: func(taskID uint64) {
				mock.ExpectQuery("SELECT (.+) FROM tasks").WillReturnError(sql.ErrNoRows)
			},
			want: want{
				err: sql.ErrNoRows,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(tt.args.taskID)

			err := uc.Delete(tt.args.ctx, tt.args.taskID)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
