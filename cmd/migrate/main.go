package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createTableRE = regexp.MustCompile(`(?i)\bCREATE\s+TABLE\b`)
	createIndexRE = regexp.MustCompile(`(?i)\bCREATE\s+(UNIQUE\s+)?INDEX\b`)
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://horizon:horizon@localhost:5433/horizon?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("database config: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database ping: %v", err)
	}
	log.Println("connected to database")

	pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		filename TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)

	migrationsDir := findMigrationsDir()
	log.Printf("migrations dir: %s", migrationsDir)

	files := listSQLFiles(migrationsDir, "V*.sql")
	for _, f := range files {
		applyIfNeeded(ctx, pool, f)
	}

	seedFiles := listSQLFiles(migrationsDir, "S*.sql")
	for _, f := range seedFiles {
		applyIfNeeded(ctx, pool, f)
	}

	fmt.Println("all migrations applied successfully")
}

func findMigrationsDir() string {
	candidates := []string{
		"services/infra/migrations",
		"../services/infra/migrations",
		"../../services/infra/migrations",
	}
	for _, d := range candidates {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			return d
		}
	}
	log.Fatalf("migrations directory not found (tried %v)", candidates)
	return ""
}

func listSQLFiles(dir, pattern string) []string {
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		log.Fatalf("list %s: %v", pattern, err)
	}
	sort.Strings(files)
	return files
}

func applyIfNeeded(ctx context.Context, pool *pgxpool.Pool, path string) {
	name := filepath.Base(path)

	var exists bool
	pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, name).Scan(&exists)
	if exists {
		log.Printf("  %s already applied, skipping", name)
		return
	}

	sql, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}

	sqlStr := makeIdempotent(string(sql))

	log.Printf("applying %s...", name)
	if _, err := pool.Exec(ctx, sqlStr); err != nil {
		log.Fatalf("execute %s: %v", name, err)
	}

	pool.Exec(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, name)
	log.Printf("  %s done", name)
}

func makeIdempotent(sql string) string {
	lines := strings.Split(sql, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			result = append(result, line)
			continue
		}
		upper := strings.ToUpper(trimmed)
		if createTableRE.MatchString(trimmed) && !strings.Contains(upper, "IF NOT EXISTS") {
			line = strings.Replace(line, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", 1)
		}
		if createIndexRE.MatchString(trimmed) && !strings.Contains(upper, "IF NOT EXISTS") {
			// Handle CREATE INDEX and CREATE UNIQUE INDEX
			line = strings.Replace(line, "CREATE UNIQUE INDEX", "CREATE UNIQUE INDEX IF NOT EXISTS", 1)
			line = strings.Replace(line, "CREATE INDEX", "CREATE INDEX IF NOT EXISTS", 1)
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
