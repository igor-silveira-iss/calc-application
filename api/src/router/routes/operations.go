package routes

import (
	"net/http"

	"github.com/igor-silveira-iss/calc-application/api/src/controllers"
)

// mathOperationsRoutes lists the calculator endpoints.
var mathOperationsRoutes = []Route{
	{
		URI:      "/addition",
		Method:   http.MethodPost,
		Function: controllers.AdditionHandler,
	},
	{
		URI:      "/subtraction",
		Method:   http.MethodPost,
		Function: controllers.SubtractionHandler,
	},
	{
		URI:      "/multiplication",
		Method:   http.MethodPost,
		Function: controllers.MultiplicationHandler,
	},
	{
		URI:      "/division",
		Method:   http.MethodPost,
		Function: controllers.DivisionHandler,
	},
	{
		URI:      "/exponentiation",
		Method:   http.MethodPost,
		Function: controllers.ExponentiationHandler,
	},
	{
		URI:      "/square-root",
		Method:   http.MethodPost,
		Function: controllers.SquareRootHandler,
	},
	{
		URI:      "/percentage",
		Method:   http.MethodPost,
		Function: controllers.PercentageHandler,
	},
}
