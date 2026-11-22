// Package response writes JSON responses.
package response

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes body as JSON with the given status code.
func WriteJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteError writes err as {"err": "<message>"}.
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteJSON(w, statusCode, struct {
		Err string `json:"err"`
	}{
		Err: err.Error(),
	})
}
