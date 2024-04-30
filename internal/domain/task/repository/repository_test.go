package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/schema"
	"github.com/henriqueassiss/advanced-golang-api/third_party/database"
	"github.com/stretchr/testify/assert"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTaskRepository_FindOne(t *testing.T) {
	db, mock := database.NewSqlxMock(t)
	r := New(db)
	defer db.Close()

	type args struct {
		ctx    context.Context
		params schema.QueryParams
	}

	type want struct {
		t   *task.Schema
		err error
	}

	type test struct {
		name string
		args
		beforeTest func()
		want
	}

	tests := []test{
		{
			name: "Success - No params",
			args: args{
				ctx:    context.TODO(),
				params: schema.QueryParams{},
			},
			beforeTest: func() {
				rows := mock.NewRows([]string{
					"id",
					"title",
					"description",
					"updated_at",
				}).AddRow(
					1,
					"Test",
					"Test",
					sql.NullTime{
						Valid: true,
					},
				)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: want{
				t: &task.Schema{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					UpdatedAt: sql.NullTime{
						Valid: true,
					},
				},
			},
		},
		{
			name: "Success - No params",
			args: args{
				ctx: context.TODO(),
				params: schema.QueryParams{
					Select: "t.id",
				},
			},
			beforeTest: func() {
				rows := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("SELECT t.id FROM tasks t").WillReturnRows(rows)
			},
			want: want{
				t: &task.Schema{
					ID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest()

			got, err := r.FindOne(tt.args.ctx, tt.args.params)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.t, got)
		})
	}
}

func TestTaskRepository_Create(t *testing.T) {
	db, mock := database.NewSqlxMock(t)
	r := New(db)
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
		beforeTest func()
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				ctx: context.TODO(),
				t: &task.Schema{
					Title:       "Test",
					Description: "Test",
				},
			},
			beforeTest: func() {
				mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlxmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest()

			err := r.Create(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {
	db, mock := database.NewSqlxMock(t)
	r := New(db)
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
		beforeTest func()
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
				},
			},
			beforeTest: func() {
				mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlxmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest()

			err := r.Update(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	db, mock := database.NewSqlxMock(t)
	r := New(db)
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
		beforeTest func()
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				ctx:    context.TODO(),
				taskID: 1,
			},
			beforeTest: func() {
				r := sqlxmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM tasks WHERE id =").WillReturnResult(r)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest()

			err := r.Delete(tt.args.ctx, tt.args.taskID)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
