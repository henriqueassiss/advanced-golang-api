package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/useCase"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/reqRes"

	"github.com/henriqueassiss/advanced-golang-api/third_party/logger"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler_FindOne(t *testing.T) {
	logger := logger.New()
	date := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

	type args struct {
		taskID string
	}

	type want struct {
		status   int
		useCase  *task.Schema
		response *reqRes.GenericResponse[*SingleTask]
		err      error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				taskID: "1",
			},
			want: want{
				status: http.StatusOK,
				useCase: &task.Schema{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					UpdatedAt: sql.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
					},
				},
				response: &reqRes.GenericResponse[*SingleTask]{
					Success: true,
					Status:  http.StatusOK,
					Data: &SingleTask{
						ID:          1,
						Title:       "Test",
						Description: "Test",
						UpdatedAt:   &date,
					},
				},
				err: nil,
			},
		},
		{
			name: "Fail - Task id not supplied",
			args: args{},
			want: want{
				status:  http.StatusBadRequest,
				useCase: &task.Schema{},
				response: &reqRes.GenericResponse[*SingleTask]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			"Fail - Non-existent task with provided task id",
			args{
				taskID: "1",
			},
			want{
				status:  http.StatusBadRequest,
				useCase: &task.Schema{},
				response: &reqRes.GenericResponse[*SingleTask]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: sql.ErrNoRows,
			},
		},
		{
			"Fail - Simulating internal error",
			args{
				taskID: "1",
			},
			want{
				status:  http.StatusInternalServerError,
				useCase: &task.Schema{},
				response: &reqRes.GenericResponse[*SingleTask]{
					Success: false,
					Status:  http.StatusInternalServerError,
					Data:    nil,
				},
				err: errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/task", nil)
			w := httptest.NewRecorder()

			if tt.args.taskID != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("taskID", tt.args.taskID)

				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			}

			uc := &useCase.TaskMock{
				FindOneFunc: func(ctx context.Context, taskID uint64) (*task.Schema, error) {
					return tt.want.useCase, tt.want.err
				},
			}

			router := chi.NewRouter()
			h := RegisterHTTPEndPoints(uc, logger, router)
			h.FindOne(w, r)

			assert.Equal(t, tt.status, w.Code)

			var got reqRes.GenericResponse[*SingleTask]
			err := json.NewDecoder(w.Body).Decode(&got)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, *tt.want.response, got)
		})
	}
}

func TestTaskHandler_Create(t *testing.T) {
	logger := logger.New()

	type args struct {
		req *Create
	}

	type want struct {
		status   int
		response *reqRes.GenericResponse[any]
		err      error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				req: &Create{
					Title:       "Test",
					Description: "Test",
				},
			},
			want: want{
				status: http.StatusOK,
				response: &reqRes.GenericResponse[any]{
					Success: true,
					Status:  http.StatusOK,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			name: "Fail - Invalid title",
			args: args{
				req: &Create{
					Description: "Test",
				},
			},
			want: want{
				status: http.StatusBadRequest,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			"Fail - Simulating internal error",
			args{
				req: &Create{
					Title:       "Test",
					Description: "Test",
				},
			},
			want{
				status: http.StatusInternalServerError,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusInternalServerError,
					Data:    nil,
				},
				err: errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			err = json.NewEncoder(&buf).Encode(tt.args.req)
			assert.Nil(t, err)

			r := httptest.NewRequest(http.MethodPost, "/api/v1/task", &buf)
			w := httptest.NewRecorder()

			uc := &useCase.TaskMock{
				CreateFunc: func(ctx context.Context, t *task.Schema) error {
					return tt.want.err
				},
			}

			router := chi.NewRouter()
			h := RegisterHTTPEndPoints(uc, logger, router)
			h.Create(w, r)

			assert.Equal(t, tt.status, w.Code)

			var got reqRes.GenericResponse[any]
			err = json.NewDecoder(w.Body).Decode(&got)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, *tt.want.response, got)
		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	logger := logger.New()

	type args struct {
		req *Update
	}

	type want struct {
		status   int
		response *reqRes.GenericResponse[any]
		err      error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				req: &Update{
					ID:          1,
					Title:       "Test",
					Description: "Test",
				},
			},
			want: want{
				status: http.StatusOK,
				response: &reqRes.GenericResponse[any]{
					Success: true,
					Status:  http.StatusOK,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			name: "Fail - Invalid id",
			args: args{
				req: &Update{
					Title:       "Test",
					Description: "Test",
				},
			},
			want: want{
				status: http.StatusBadRequest,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			name: "Fail - Invalid title",
			args: args{
				req: &Update{
					ID:          1,
					Description: "Test",
				},
			},
			want: want{
				status: http.StatusBadRequest,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			"Fail - Simulating internal error",
			args{
				req: &Update{
					ID:          1,
					Title:       "Test",
					Description: "Test",
				},
			},
			want{
				status: http.StatusInternalServerError,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusInternalServerError,
					Data:    nil,
				},
				err: errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			err = json.NewEncoder(&buf).Encode(tt.args.req)
			assert.Nil(t, err)

			r := httptest.NewRequest(http.MethodPut, "/api/v1/task", &buf)
			w := httptest.NewRecorder()

			uc := &useCase.TaskMock{
				UpdateFunc: func(ctx context.Context, t *task.Schema) error {
					return tt.want.err
				},
			}

			router := chi.NewRouter()
			h := RegisterHTTPEndPoints(uc, logger, router)
			h.Update(w, r)

			assert.Equal(t, tt.status, w.Code)

			var got reqRes.GenericResponse[any]
			err = json.NewDecoder(w.Body).Decode(&got)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, *tt.want.response, got)
		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	logger := logger.New()

	type args struct {
		taskID string
	}

	type want struct {
		status   int
		response *reqRes.GenericResponse[any]
		err      error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "Success",
			args: args{
				taskID: "1",
			},
			want: want{
				status: http.StatusOK,
				response: &reqRes.GenericResponse[any]{
					Success: true,
					Status:  http.StatusOK,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			name: "Fail - Task id not supplied",
			args: args{},
			want: want{
				status: http.StatusBadRequest,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: nil,
			},
		},
		{
			name: "Fail - Non-existent task with provided task id",
			args: args{
				taskID: "1",
			},
			want: want{
				status: http.StatusBadRequest,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusBadRequest,
					Data:    nil,
				},
				err: sql.ErrNoRows,
			},
		},
		{
			name: "Fail - Simulating internal error",
			args: args{
				taskID: "1",
			},
			want: want{
				status: http.StatusInternalServerError,
				response: &reqRes.GenericResponse[any]{
					Success: false,
					Status:  http.StatusInternalServerError,
					Data:    nil,
				},
				err: errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/task", nil)
			w := httptest.NewRecorder()

			if tt.args.taskID != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("taskID", tt.args.taskID)

				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			}

			uc := &useCase.TaskMock{
				DeleteFunc: func(ctx context.Context, taskID uint64) error {
					return tt.want.err
				},
			}

			router := chi.NewRouter()
			h := RegisterHTTPEndPoints(uc, logger, router)
			h.Delete(w, r)

			assert.Equal(t, tt.status, w.Code)

			var got reqRes.GenericResponse[any]
			err := json.NewDecoder(w.Body).Decode(&got)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, *tt.want.response, got)
		})
	}
}
