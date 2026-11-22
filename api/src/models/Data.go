// Package models defines the input and output types of the operations.
package models

// Input holds the two numbers of an operation.
type Input struct {
	Number1 float64
	Number2 float64
}

// InputSquareRoot holds the number of a square root.
type InputSquareRoot struct {
	Number float64
}

// Output holds the result of an operation, sent as {"result": <number>}.
type Output struct {
	Result float64 `json:"result"`
}
