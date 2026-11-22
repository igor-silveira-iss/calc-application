package controllers

import (
	"errors"
	"net/http"

	"github.com/igor-silveira-iss/calc-application/api/src/models"
	"github.com/igor-silveira-iss/calc-application/api/src/response"
)

// AdditionHandler handles POST /addition.
func AdditionHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Addition)
}

// SubtractionHandler handles POST /subtraction.
func SubtractionHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Subtraction)
}

// MultiplicationHandler handles POST /multiplication.
func MultiplicationHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Multiplication)
}

// DivisionHandler handles POST /division.
func DivisionHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Division)
}

// ExponentiationHandler handles POST /exponentiation.
func ExponentiationHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Exponentiation)
}

// PercentageHandler handles POST /percentage.
func PercentageHandler(w http.ResponseWriter, r *http.Request) {
	handleOperands(w, r, Percentage)
}

// SquareRootHandler handles POST /square-root.
func SquareRootHandler(w http.ResponseWriter, r *http.Request) {
	var body numberRequest
	if !decodeBody(w, r, &body) {
		return
	}
	if body.Number == nil {
		response.WriteError(w, http.StatusBadRequest, errors.New("number is required"))
		return
	}
	output, err := SquareRoot(models.InputSquareRoot{Number: *body.Number})
	writeResult(w, output, err)
}

// operandsRequest is the body of a two-number operation.
// Pointers tell a missing field apart from zero.
type operandsRequest struct {
	Number1 *float64 `json:"number1"`
	Number2 *float64 `json:"number2"`
}

// numberRequest is the body of the square root operation.
type numberRequest struct {
	Number *float64 `json:"number"`
}

// handleOperands decodes number1 and number2 and writes the result of operation.
func handleOperands(w http.ResponseWriter, r *http.Request, operation func(models.Input) (models.Output, error)) {
	var body operandsRequest
	if !decodeBody(w, r, &body) {
		return
	}
	switch {
	case body.Number1 == nil:
		response.WriteError(w, http.StatusBadRequest, errors.New("number1 is required"))
		return
	case body.Number2 == nil:
		response.WriteError(w, http.StatusBadRequest, errors.New("number2 is required"))
		return
	}
	output, err := operation(models.Input{Number1: *body.Number1, Number2: *body.Number2})
	writeResult(w, output, err)
}

// writeResult writes output with 200, or err with 422.
func writeResult(w http.ResponseWriter, output models.Output, err error) {
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, output)
}
