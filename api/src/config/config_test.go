package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/igor-silveira-iss/calc-application/api/src/config"
)

// useEnvFile runs the test in a temporary directory with a .env file holding
// content, and clears the variables so the file values are used.
func useEnvFile(t *testing.T, content string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	for _, key := range []string{"API_PORT", "ALLOWED_ORIGIN"} {
		t.Setenv(key, "") // restores the original value after the test
		os.Unsetenv(key)
	}
}

func TestInitEnv(t *testing.T) {
	tests := []struct {
		name       string
		envFile    string
		noEnvFile  bool
		wantPort   int
		wantOrigin string
	}{
		{
			name:       "values from .env",
			envFile:    "API_PORT=9000\nALLOWED_ORIGIN=http://example.com\n",
			wantPort:   9000,
			wantOrigin: "http://example.com",
		},
		{
			name:       "defaults when empty",
			envFile:    "",
			wantPort:   8080,
			wantOrigin: "http://localhost:5173",
		},
		{
			name:       "defaults without .env file",
			envFile:    "",
			noEnvFile:  true,
			wantPort:   8080,
			wantOrigin: "http://localhost:5173",
		},
		{
			name:       "default port when invalid",
			envFile:    "API_PORT=abc\n",
			wantPort:   8080,
			wantOrigin: "http://localhost:5173",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useEnvFile(t, tt.envFile)
			if tt.noEnvFile {
				os.Remove(".env")
			}
			config.InitEnv()

			if config.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", config.Port, tt.wantPort)
			}
			if config.AllowedOrigin != tt.wantOrigin {
				t.Errorf("AllowedOrigin = %q, want %q", config.AllowedOrigin, tt.wantOrigin)
			}
		})
	}
}
