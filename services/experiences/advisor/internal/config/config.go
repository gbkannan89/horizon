package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() *Config {
	p := os.Getenv("PORT")
	if p == "" { p = "8097" }
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable"
	}
	return &Config{Port: p, DatabaseURL: db}
}

func (c *Config) Addr() string { return ":" + c.Port }
