package config

import (
	"os"
	"strconv"
)

// Config holds runtime configuration.
type Config struct {
	HTTPPort int
}

// Load reads the application configuration from environment variables.
func Load() Config {
	port := 8080
	if raw := os.Getenv("HTTP_PORT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			port = parsed
		}
	}

	return Config{
		HTTPPort: port,
	}
}
