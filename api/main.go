// Command api starts the calculator HTTP API.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igor-silveira-iss/calc-application/api/src/config"
	"github.com/igor-silveira-iss/calc-application/api/src/middleware"
	"github.com/igor-silveira-iss/calc-application/api/src/router"
)

// main loads the config and starts the HTTP server.
func main() {
	config.InitEnv()
	r := router.StartRouter()

	fmt.Printf("Listening on port %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), middleware.CORS(config.AllowedOrigin, r)))
}
