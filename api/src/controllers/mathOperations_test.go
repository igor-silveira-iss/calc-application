package controllers_test

import (
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/controllers"
	"github.com/igor-silveira-iss/calc-application/api/src/models"
)

const (
	errOutOfRange     = "result exceeds the supported numeric range"
	errDivisionByZero = "division by zero is not allowed"
)

type operationTest struct {
	name    string
	input   models.Input
	want    float64
	wantErr string
}

// runOperationTests checks each test case against operation.
func runOperationTests(t *testing.T, operation func(models.Input) (models.Output, error), tests []operationTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := operation(tt.input)
			checkResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

// checkResult compares an operation result with the expected value or error.
func checkResult(t *testing.T, got models.Output, err error, want float64, wantErr string) {
	t.Helper()
	if wantErr != "" {
		if err == nil || err.Error() != wantErr {
			t.Fatalf("error = %v, want %q", err, wantErr)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Result != want {
		t.Errorf("result = %v, want %v", got.Result, want)
	}
}

func TestAddition(t *testing.T) {
	runOperationTests(t, controllers.Addition, []operationTest{
		{name: "positive numbers", input: models.Input{Number1: 2, Number2: 3}, want: 5},
		{name: "negative numbers", input: models.Input{Number1: -2, Number2: -3}, want: -5},
		{name: "decimals", input: models.Input{Number1: 1.5, Number2: 2.25}, want: 3.75},
		{name: "zeros", input: models.Input{Number1: 0, Number2: 0}, want: 0},
		{name: "overflow", input: models.Input{Number1: 1e308, Number2: 1e308}, wantErr: errOutOfRange},
	})
}

func TestSubtraction(t *testing.T) {
	runOperationTests(t, controllers.Subtraction, []operationTest{
		{name: "positive result", input: models.Input{Number1: 9, Number2: 3}, want: 6},
		{name: "negative result", input: models.Input{Number1: 3, Number2: 9}, want: -6},
		{name: "decimals", input: models.Input{Number1: 5.5, Number2: 2.25}, want: 3.25},
		{name: "overflow", input: models.Input{Number1: -1e308, Number2: 1e308}, wantErr: errOutOfRange},
	})
}

func TestMultiplication(t *testing.T) {
	runOperationTests(t, controllers.Multiplication, []operationTest{
		{name: "positive numbers", input: models.Input{Number1: 9, Number2: 3}, want: 27},
		{name: "negative number", input: models.Input{Number1: -4, Number2: 2.5}, want: -10},
		{name: "by zero", input: models.Input{Number1: 7, Number2: 0}, want: 0},
		{name: "overflow", input: models.Input{Number1: 1e200, Number2: 1e200}, wantErr: errOutOfRange},
	})
}

func TestDivision(t *testing.T) {
	runOperationTests(t, controllers.Division, []operationTest{
		{name: "exact division", input: models.Input{Number1: 9, Number2: 3}, want: 3},
		{name: "decimal result", input: models.Input{Number1: 1, Number2: 4}, want: 0.25},
		{name: "negative divisor", input: models.Input{Number1: 10, Number2: -2}, want: -5},
		{name: "zero dividend", input: models.Input{Number1: 0, Number2: 5}, want: 0},
		{name: "division by zero", input: models.Input{Number1: 1, Number2: 0}, wantErr: errDivisionByZero},
		{name: "zero by zero", input: models.Input{Number1: 0, Number2: 0}, wantErr: errDivisionByZero},
		{name: "overflow", input: models.Input{Number1: 1e308, Number2: 1e-10}, wantErr: errOutOfRange},
	})
}

func TestExponentiation(t *testing.T) {
	runOperationTests(t, controllers.Exponentiation, []operationTest{
		{name: "integer power", input: models.Input{Number1: 2, Number2: 10}, want: 1024},
		{name: "negative base and odd power", input: models.Input{Number1: -2, Number2: 3}, want: -8},
		{name: "negative base and negative power", input: models.Input{Number1: -2, Number2: -2}, want: 0.25},
		{name: "negative power", input: models.Input{Number1: 2, Number2: -1}, want: 0.5},
		{name: "fractional power", input: models.Input{Number1: 4, Number2: 0.5}, want: 2},
		{name: "zero to the zero", input: models.Input{Number1: 0, Number2: 0}, want: 1},
		{name: "zero to a negative power", input: models.Input{Number1: 0, Number2: -1},
			wantErr: "zero cannot be raised to a negative power"},
		{name: "negative base and fractional power", input: models.Input{Number1: -8, Number2: 0.5},
			wantErr: "negative number cannot be raised to a fractional power"},
		{name: "overflow", input: models.Input{Number1: 10, Number2: 400}, wantErr: errOutOfRange},
	})
}

func TestPercentage(t *testing.T) {
	runOperationTests(t, controllers.Percentage, []operationTest{
		{name: "part of a whole", input: models.Input{Number1: 25, Number2: 200}, want: 12.5},
		{name: "more than the whole", input: models.Input{Number1: 300, Number2: 200}, want: 150},
		{name: "negative part", input: models.Input{Number1: -50, Number2: 200}, want: -25},
		{name: "zero part", input: models.Input{Number1: 0, Number2: 200}, want: 0},
		{name: "zero whole", input: models.Input{Number1: 5, Number2: 0}, wantErr: errDivisionByZero},
		{name: "zero of zero", input: models.Input{Number1: 0, Number2: 0}, wantErr: errDivisionByZero},
		{name: "overflow", input: models.Input{Number1: 1e308, Number2: 1e-10}, wantErr: errOutOfRange},
	})
}

func TestSquareRoot(t *testing.T) {
	tests := []struct {
		name    string
		number  float64
		want    float64
		wantErr string
	}{
		{name: "perfect square", number: 16, want: 4},
		{name: "decimal", number: 2.25, want: 1.5},
		{name: "zero", number: 0, want: 0},
		{name: "negative number", number: -4, wantErr: "square root of a negative number is not allowed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := controllers.SquareRoot(models.InputSquareRoot{Number: tt.number})
			checkResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}
