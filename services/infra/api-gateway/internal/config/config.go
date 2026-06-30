package config

import "os"

type Config struct {
	Port          string
	AllowedOrigins []string
	LogLevel      string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		AllowedOrigins: []string{"*"},
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}
