package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/useCase"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/errorMsg"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/reqRes"
)

type ITask struct {
	useCase useCase.ITask
	logger  *slog.Logger
}

func NewHandler(useCase useCase.ITask, logger *slog.Logger) *ITask {
	return &ITask{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *ITask) FindOne(w http.ResponseWriter, r *http.Request) {
	taskID, err := reqRes.UInt64Param(r, "taskID", false)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusBadRequest, errorMsg.ErrInvalidRequestData, taskID)
		return
	}

	t, err := h.useCase.FindOne(r.Context(), taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			reqRes.Error(h.logger, w, http.StatusBadRequest, err, taskID)
		} else {
			reqRes.Error(h.logger, w, http.StatusInternalServerError, err, taskID)
		}

		return
	}

	res := SingleTask{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
	}

	if t.UpdatedAt.Valid {
		res.UpdatedAt = &t.UpdatedAt.Time
	}

	reqRes.Json(w, http.StatusOK, res)
}

func (h *ITask) Create(w http.ResponseWriter, r *http.Request) {
	var req Create
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusBadRequest, err, req)
		return
	}

	if req.Title == "" {
		reqRes.Error(h.logger, w, http.StatusBadRequest, nil, req)
		return
	}

	t := task.Schema{
		Title:       req.Title,
		Description: req.Description,
	}

	err = h.useCase.Create(r.Context(), &t)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusInternalServerError, err, t)
		return
	}

	reqRes.Json(w, http.StatusOK, nil)
}

func (h *ITask) Update(w http.ResponseWriter, r *http.Request) {
	var req Update
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusBadRequest, err, req)
		return
	}

	if req.ID == 0 ||
		req.Title == "" {
		reqRes.Error(h.logger, w, http.StatusBadRequest, nil, req)
		return
	}

	t := task.Schema{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	err = h.useCase.Update(r.Context(), &t)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusInternalServerError, err, t)
		return
	}

	reqRes.Json(w, http.StatusOK, nil)
}

func (h *ITask) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := reqRes.UInt64Param(r, "taskID", false)
	if err != nil {
		reqRes.Error(h.logger, w, http.StatusBadRequest, err, taskID)
		return
	}

	err = h.useCase.Delete(r.Context(), taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			reqRes.Error(h.logger, w, http.StatusBadRequest, err, taskID)
		} else {
			reqRes.Error(h.logger, w, http.StatusInternalServerError, err, taskID)
		}

		return
	}

	reqRes.Json(w, http.StatusOK, nil)
}
