package config

// Config holds the application layer configuration.
type Config struct {
	JWTSecret         string `env:"JWT_SECRET,required"`
	DatabaseURL       string `env:"DATABASE_URL,required"`
	RedisAddr         string `env:"REDIS_ADDR,required"`
	NATSURL           string `env:"NATS_URL,required"`
	LogLevel          string `env:"LOG_LEVEL" default:"info"`
	RateLimitPerMin   int    `env:"RATE_LIMIT_PER_MIN" default:"1000"`
	IdempotencyTTL    int    `env:"IDEMPOTENCY_TTL_HOURS" default:"24"`
}

// DefaultConfig returns configuration with sensible defaults for development.
func DefaultConfig() Config {
	return Config{
		JWTSecret:       "dev-secret-key-do-not-use-in-production",
		DatabaseURL:     "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable",
		RedisAddr:       "localhost:6379",
		NATSURL:         "nats://localhost:4222",
		LogLevel:        "debug",
		RateLimitPerMin: 1000,
		IdempotencyTTL:  24,
	}
}
