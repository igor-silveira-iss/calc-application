// Package controllers implements the calculator operations.
package controllers

import (
	"errors"
	"math"

	"github.com/igor-silveira-iss/calc-application/api/src/models"
)

// Addition returns Number1 + Number2.
func Addition(data models.Input) (models.Output, error) {
	result := data.Number1 + data.Number2
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// Subtraction returns Number1 - Number2.
func Subtraction(data models.Input) (models.Output, error) {
	result := data.Number1 - data.Number2
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// Multiplication returns Number1 * Number2.
func Multiplication(data models.Input) (models.Output, error) {
	result := data.Number1 * data.Number2
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// Division returns Number1 / Number2. It fails if Number2 is zero.
func Division(data models.Input) (models.Output, error) {
	if data.Number2 == 0 {
		return models.Output{}, errors.New("division by zero is not allowed")
	}
	result := data.Number1 / data.Number2
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// Exponentiation returns Number1 raised to Number2.
// It fails if zero is raised to a negative power or a negative number is
// raised to a fractional power.
func Exponentiation(data models.Input) (models.Output, error) {
	if data.Number1 == 0 && data.Number2 < 0 {
		return models.Output{}, errors.New("zero cannot be raised to a negative power")
	}
	if data.Number1 < 0 && data.Number2 != math.Trunc(data.Number2) {
		return models.Output{}, errors.New("negative number cannot be raised to a fractional power")
	}
	result := math.Pow(data.Number1, data.Number2)
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// Percentage returns how many percent Number1 is of Number2.
func Percentage(data models.Input) (models.Output, error) {
	if data.Number2 == 0 {
		return models.Output{}, errors.New("division by zero is not allowed")
	}
	result := (data.Number1 / data.Number2) * 100
	if math.IsInf(result, 0) {
		return models.Output{}, errors.New("result exceeds the supported numeric range")
	}
	return models.Output{Result: result}, nil
}

// SquareRoot returns the square root of Number.
func SquareRoot(data models.InputSquareRoot) (models.Output, error) {
	if data.Number < 0 {
		return models.Output{}, errors.New("square root of a negative number is not allowed")
	}
	return models.Output{Result: math.Sqrt(data.Number)}, nil
}
