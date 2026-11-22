package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/response"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	response.WriteJSON(rec, http.StatusCreated, map[string]int{"value": 1})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"value":1}` {
		t.Errorf("body = %s, want {\"value\":1}", got)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	response.WriteError(rec, http.StatusBadRequest, errors.New("something went wrong"))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"err":"something went wrong"}` {
		t.Errorf("body = %s, want {\"err\":\"something went wrong\"}", got)
	}
}
