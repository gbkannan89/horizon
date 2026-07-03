package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Provider    string        `json:"provider"`
	Port        string        `json:"port"`
	MaxHistory  int           `json:"max_history"`
	AIEnabled   bool          `json:"ai_enabled"`
	OllamaURL   string        `json:"ollama_url"`
	OllamaModel string        `json:"ollama_model"`
	Timeout     time.Duration `json:"timeout"`
	Retries     int           `json:"retries"`
	CBThreshold int           `json:"cb_threshold"`
	CBTimeout   time.Duration `json:"cb_timeout"`
}

func Load() *Config {
	p := env("AI_PROVIDER", "mock")
	port := env("PORT", "8200")
	enabled := env("AI_ENABLED", "true") != "false"
	ollamaURL := env("OLLAMA_URL", "http://localhost:11434")
	ollamaModel := env("OLLAMA_MODEL", "llama3")
	timeoutSec := intEnv("AI_TIMEOUT_SECONDS", 30)
	retries := intEnv("AI_RETRIES", 2)
	cbThreshold := intEnv("AI_CB_THRESHOLD", 3)
	cbTimeoutSec := intEnv("AI_CB_TIMEOUT_SECONDS", 60)

	return &Config{
		Provider: p, Port: port, MaxHistory: 50,
		AIEnabled: enabled, OllamaURL: ollamaURL,
		OllamaModel: ollamaModel, Timeout: time.Duration(timeoutSec) * time.Second,
		Retries: retries, CBThreshold: cbThreshold, CBTimeout: time.Duration(cbTimeoutSec) * time.Second,
	}
}

func (c *Config) Addr() string { return ":" + c.Port }

func env(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func intEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil { return n }
	}
	return def
}
