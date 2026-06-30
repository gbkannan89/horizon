package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps a PostgreSQL connection pool.
type DB struct {
	pool *sql.DB
}

// NewDB creates a new database connection.
func NewDB(databaseURL string) (*DB, error) {
	pool, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(5 * time.Minute)
	return &DB{pool: pool}, nil
}

// Ping checks database connectivity.
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.PingContext(ctx)
}

// Close closes the connection pool.
func (db *DB) Close() error {
	return db.pool.Close()
}

// Pool returns the underlying connection pool.
func (db *DB) Pool() *sql.DB { return db.pool }

// HealthCheck verifies database reachability.
func (db *DB) HealthCheck(ctx context.Context) map[string]interface{} {
	err := db.Ping(ctx)
	status := "healthy"
	if err != nil {
		status = "unhealthy"
	}
	return map[string]interface{}{
		"status":     status,
		"error":      err,
		"open_conns": db.pool.Stats().OpenConnections,
	}
}
