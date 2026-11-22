package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/middleware"
)

func TestLogger(t *testing.T) {
	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	called := false
	handler := middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/addition", nil))

	if !called {
		t.Error("next handler was not called")
	}
	if !strings.Contains(logs.String(), "POST /addition") {
		t.Errorf("log = %q, want it to contain %q", logs.String(), "POST /addition")
	}
}

func TestCORS(t *testing.T) {
	const origin = "http://localhost:5173"

	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantNext   bool
	}{
		{name: "preflight request", method: http.MethodOptions, wantStatus: http.StatusNoContent, wantNext: false},
		{name: "regular request", method: http.MethodPost, wantStatus: http.StatusOK, wantNext: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})
			rec := httptest.NewRecorder()
			middleware.CORS(origin, next).ServeHTTP(rec, httptest.NewRequest(tt.method, "/addition", nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called != tt.wantNext {
				t.Errorf("next called = %v, want %v", called, tt.wantNext)
			}
			wantHeaders := map[string]string{
				"Access-Control-Allow-Origin":  origin,
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
				"Vary":                         "Origin",
			}
			for header, want := range wantHeaders {
				if got := rec.Header().Get(header); got != want {
					t.Errorf("%s = %q, want %q", header, got, want)
				}
			}
		})
	}
}
