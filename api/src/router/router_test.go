package router_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/router"
)

func TestStartRouter(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{"addition", http.MethodPost, "/addition", `{"number1":2,"number2":3}`, 200, `{"result":5}`},
		{"subtraction", http.MethodPost, "/subtraction", `{"number1":2,"number2":3}`, 200, `{"result":-1}`},
		{"multiplication", http.MethodPost, "/multiplication", `{"number1":2,"number2":3}`, 200, `{"result":6}`},
		{"division", http.MethodPost, "/division", `{"number1":3,"number2":2}`, 200, `{"result":1.5}`},
		{"exponentiation", http.MethodPost, "/exponentiation", `{"number1":2,"number2":3}`, 200, `{"result":8}`},
		{"percentage", http.MethodPost, "/percentage", `{"number1":1,"number2":4}`, 200, `{"result":25}`},
		{"square root", http.MethodPost, "/square-root", `{"number":9}`, 200, `{"result":3}`},
		{"unknown route", http.MethodPost, "/unknown", ``, 404, `{"err":"route not found"}`},
		{"wrong method", http.MethodGet, "/addition", ``, 405, `{"err":"method not allowed"}`},
	}

	r := router.StartRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body)))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}
