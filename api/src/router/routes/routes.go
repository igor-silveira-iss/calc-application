// Package routes defines the API routes.
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/igor-silveira-iss/calc-application/api/src/middleware"
)

// Route describes an API endpoint.
type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
}

// Configure registers all routes on the router.
func Configure(r *mux.Router) *mux.Router {
	routes := mathOperationsRoutes
	for _, route := range routes {
		r.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
	}
	return r
}
