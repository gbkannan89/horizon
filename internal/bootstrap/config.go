package bootstrap

import "os"

// Config holds all application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	AppEnv      string
	LogLevel    string
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable"),
		AppEnv:      env("APP_ENV", "development"),
		LogLevel:    env("LOG_LEVEL", "info"),
	}
}

// Addr returns the address string for the HTTP server.
func (c *Config) Addr() string { return ":" + c.Port }

func env(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}
