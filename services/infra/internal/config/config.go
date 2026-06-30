package config

import "os"

// Config holds all infrastructure configuration.
type Config struct {
	DatabaseURL string `env:"DATABASE_URL" default:"postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable"`
	RedisAddr   string `env:"REDIS_ADDR" default:"localhost:6379"`
	NATSURL     string `env:"NATS_URL" default:"nats://localhost:4222"`
	StorageURL  string `env:"STORAGE_ENDPOINT" default:"http://localhost:9000"`
	LogLevel    string `env:"LOG_LEVEL" default:"info"`
}

// Load loads configuration from environment variables with defaults.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		StorageURL:  getEnv("STORAGE_ENDPOINT", "http://localhost:9000"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}
