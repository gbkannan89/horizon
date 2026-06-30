package config

import "os"

type Config struct {
	Port string
}

func Load() *Config {
	p := os.Getenv("PORT")
	if p == "" { p = "8099" }
	return &Config{Port: p}
}

func (c *Config) Addr() string { return ":" + c.Port }
