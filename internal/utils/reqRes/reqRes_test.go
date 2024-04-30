package reqRes

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUInt64Param(t *testing.T) {
	router := chi.NewRouter()
	w := httptest.NewRecorder()

	validValue := 100
	invalidValue := "abc"

	expectedValue, expectedErr := uint64(100), error(nil)
	router.Get("/test/{param}", func(w http.ResponseWriter, r *http.Request) {
		value, err := UInt64Param(r, "param", false)
		if value != expectedValue || err != expectedErr {
			t.Errorf("got: value = %d and err = %v | expected: value = %d and err = %v", value, err, expectedValue, expectedErr)
		}
	})

	r := httptest.NewRequest("GET", fmt.Sprintf("/test/%d", validValue), nil)
	router.ServeHTTP(w, r)

	expectedValue, expectedErr = uint64(0), errors.New(`strconv.ParseInt: parsing "abc": invalid syntax`)
	router.Get("/test/{param}", func(w http.ResponseWriter, r *http.Request) {
		value, err := UInt64Param(r, "param", false)
		if value != expectedValue || err.Error() != expectedErr.Error() {
			t.Errorf("got: value = %d and err = %v | expected: value = %d and err = %v", value, err, expectedValue, expectedErr)
		}
	})

	r = httptest.NewRequest("GET", fmt.Sprintf("/test/%s", invalidValue), nil)
	router.ServeHTTP(w, r)
}

func TestUInt64Query(t *testing.T) {
	r := httptest.NewRequest("GET", "/test?param=100", nil)

	expectedValue, expectedErr := uint64(100), error(nil)
	value, err := UInt64Query(r, "param", false)
	if value != expectedValue || err != expectedErr {
		t.Errorf("got: value = %d and err = %v | expected: value = %d and err = %v", value, err, expectedValue, expectedErr)
	}
}
