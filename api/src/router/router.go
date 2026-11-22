// Package router creates the HTTP router.
package router

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/igor-silveira-iss/calc-application/api/src/response"
	"github.com/igor-silveira-iss/calc-application/api/src/router/routes"
)

// StartRouter returns a router with all routes registered.
func StartRouter() *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response.WriteError(w, http.StatusNotFound, errors.New("route not found"))
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response.WriteError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	})
	return routes.Configure(r)
}
