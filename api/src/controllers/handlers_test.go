package controllers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/controllers"
)

// serve sends a POST request with body to handler and returns the recorded response.
func serve(handler http.HandlerFunc, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

// checkResponse compares the status code and JSON body of a response.
func checkResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantBody string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Errorf("status = %d, want %d", rec.Code, wantStatus)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != wantBody {
		t.Errorf("body = %s, want %s", got, wantBody)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
}

func TestOperationHandlers(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		body       string
		wantStatus int
		wantBody   string
	}{
		{"addition", controllers.AdditionHandler, `{"number1":2,"number2":3}`, 200, `{"result":5}`},
		{"subtraction", controllers.SubtractionHandler, `{"number1":9,"number2":3}`, 200, `{"result":6}`},
		{"multiplication", controllers.MultiplicationHandler, `{"number1":9,"number2":3}`, 200, `{"result":27}`},
		{"division", controllers.DivisionHandler, `{"number1":9,"number2":3}`, 200, `{"result":3}`},
		{"exponentiation", controllers.ExponentiationHandler, `{"number1":2,"number2":10}`, 200, `{"result":1024}`},
		{"percentage", controllers.PercentageHandler, `{"number1":25,"number2":200}`, 200, `{"result":12.5}`},
		{"square root", controllers.SquareRootHandler, `{"number":16}`, 200, `{"result":4}`},
		{"zeros are valid input", controllers.AdditionHandler, `{"number1":0,"number2":0}`, 200, `{"result":0}`},

		{"division by zero", controllers.DivisionHandler, `{"number1":1,"number2":0}`,
			422, `{"err":"division by zero is not allowed"}`},
		{"percentage of zero", controllers.PercentageHandler, `{"number1":0,"number2":0}`,
			422, `{"err":"division by zero is not allowed"}`},
		{"negative square root", controllers.SquareRootHandler, `{"number":-4}`,
			422, `{"err":"square root of a negative number is not allowed"}`},
		{"negative base and fractional power", controllers.ExponentiationHandler, `{"number1":-8,"number2":0.5}`,
			422, `{"err":"negative number cannot be raised to a fractional power"}`},
		{"overflow", controllers.AdditionHandler, `{"number1":1e308,"number2":1e308}`,
			422, `{"err":"result exceeds the supported numeric range"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(tt.handler, strings.NewReader(tt.body))
			checkResponse(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestHandlersRejectInvalidBody(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		body     string
		wantBody string
	}{
		{"empty body", controllers.AdditionHandler, ``, `{"err":"request body is required"}`},
		{"empty object", controllers.AdditionHandler, `{}`, `{"err":"number1 is required"}`},
		{"missing number2", controllers.AdditionHandler, `{"number1":2}`, `{"err":"number2 is required"}`},
		{"null body", controllers.AdditionHandler, `null`, `{"err":"number1 is required"}`},
		{"null field", controllers.AdditionHandler, `{"number1":null,"number2":2}`, `{"err":"number1 is required"}`},
		{"missing number", controllers.SquareRootHandler, `{}`, `{"err":"number is required"}`},
		{"unknown field", controllers.AdditionHandler, `{"num1":5,"num2":3}`, `{"err":"unknown field \"num1\""}`},
		{"string instead of number", controllers.AdditionHandler, `{"number1":"abc","number2":1}`,
			`{"err":"number1 must be a valid number"}`},
		{"number out of range", controllers.AdditionHandler, `{"number1":1e400,"number2":1}`,
			`{"err":"number1 must be a valid number"}`},
		{"invalid square root number", controllers.SquareRootHandler, `{"number":true}`,
			`{"err":"number must be a valid number"}`},
		{"malformed JSON", controllers.AdditionHandler, `{"number1":1,`, `{"err":"request body must be valid JSON"}`},
		{"invalid JSON syntax", controllers.AdditionHandler, `{number1}`, `{"err":"request body must be valid JSON"}`},
		{"array instead of object", controllers.AdditionHandler, `[1,2]`, `{"err":"request body must be a JSON object"}`},
		{"two JSON objects", controllers.AdditionHandler, `{"number1":1,"number2":2}{"x":1}`,
			`{"err":"request body must contain a single JSON object"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(tt.handler, strings.NewReader(tt.body))
			checkResponse(t, rec, http.StatusBadRequest, tt.wantBody)
		})
	}
}

func TestHandlersRejectLargeBody(t *testing.T) {
	body := `{"number1":1,"number2":2` + strings.Repeat(" ", 2048) + `}`
	rec := serve(controllers.AdditionHandler, strings.NewReader(body))
	checkResponse(t, rec, http.StatusRequestEntityTooLarge, `{"err":"request body must not exceed 1024 bytes"}`)
}

// failingReader simulates a request body that cannot be read.
type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("connection reset")
}

func TestHandlersRejectUnreadableBody(t *testing.T) {
	rec := serve(controllers.AdditionHandler, failingReader{})
	checkResponse(t, rec, http.StatusBadRequest, `{"err":"connection reset"}`)
}
