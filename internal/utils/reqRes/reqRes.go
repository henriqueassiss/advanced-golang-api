package reqRes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mssola/useragent"
)

type GenericResponse[T any] struct {
	Success bool `json:"success"`
	Status  int  `json:"status"`
	Data    T    `json:"data"`
}

func respond(w http.ResponseWriter, statusCode int, isSuccess bool, payload any) {
	res := GenericResponse[any]{
		Success: isSuccess,
		Status:  statusCode,
		Data:    payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}

func Json(w http.ResponseWriter, statusCode int, payload any) {
	respond(w, statusCode, true, payload)
}

func Error(logger *slog.Logger, w http.ResponseWriter, statusCode int, err error, errData any) {
	respond(w, statusCode, false, nil)
	logger.Error(err.Error(), errData)
}

func GetRequestDevice(reqUserAgent string) (device string) {
	result := useragent.New(reqUserAgent)

	browser, _ := result.Browser()
	if browser == "" {
		browser = "Unknown"
	}

	isMobile := result.Mobile()
	model := result.Model()
	if model == "" {
		model = "Unknown"
	}

	platform := result.Platform()
	if platform == "" {
		platform = "Unknown"
	}

	os := result.OS()
	if os == "" {
		os = "Unknown"
	}

	device = fmt.Sprintf("Browser: %s | IsMobile: %v | Model: %s | Platform: %s | OS: %s", browser, isMobile, model, platform, os)

	return device
}

func UInt64Param(r *http.Request, param string, acceptZero bool) (uint64, error) {
	val, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, err
	}

	if !acceptZero && val == 0 {
		return 0, errors.New("run-time: invalid zero query param")
	}

	return uint64(val), err
}

func UInt64Query(r *http.Request, param string, acceptZero bool) (uint64, error) {
	val, err := strconv.ParseInt(r.URL.Query().Get(param), 10, 64)
	if err != nil {
		return 0, err
	}

	if !acceptZero && val == 0 {
		return 0, errors.New("run-time: invalid zero query param")
	}

	return uint64(val), err
}
