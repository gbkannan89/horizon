package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/rules/register"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8120" }
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { dsn = "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable" }
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil { log.Fatalf("database: %v", err) }
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"live"}`))
	})
	register.RegisterRoutes(mux, pool)
	log.Printf("Rules domain service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
