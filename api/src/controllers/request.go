package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/igor-silveira-iss/calc-application/api/src/response"
)

// maxBodyBytes is the largest request body accepted.
const maxBodyBytes = 1024

// decodeBody reads the JSON body into dst. On failure it writes the error
// response and returns false.
func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err == nil && decoder.Decode(&struct{}{}) != io.EOF {
		err = errors.New("request body must contain a single JSON object")
	}
	if err == nil {
		return true
	}

	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		response.WriteError(w, http.StatusRequestEntityTooLarge,
			fmt.Errorf("request body must not exceed %d bytes", maxBodyBytes))
		return false
	}
	response.WriteError(w, http.StatusBadRequest, bodyError(err))
	return false
}

// bodyError turns a JSON decoding error into a message for the client.
func bodyError(err error) error {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.Is(err, io.EOF):
		return errors.New("request body is required")
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("request body must be valid JSON")
	case errors.As(err, &typeErr) && typeErr.Field != "":
		return fmt.Errorf("%s must be a valid number", typeErr.Field)
	case errors.As(err, &typeErr):
		return errors.New("request body must be a JSON object")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		return errors.New(strings.TrimPrefix(err.Error(), "json: "))
	default:
		return err
	}
}
