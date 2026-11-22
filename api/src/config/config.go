// Package config loads the application settings.
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Port is the port the server listens on.
var Port = 0

// AllowedOrigin is the frontend URL allowed by CORS.
var AllowedOrigin = ""

// InitEnv loads the .env file, if any, and sets Port (default 8080) and
// AllowedOrigin (default http://localhost:5173).
func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	var err error
	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 8080
	}
	AllowedOrigin = os.Getenv("ALLOWED_ORIGIN")
	if AllowedOrigin == "" {
		AllowedOrigin = "http://localhost:5173"
	}
}
