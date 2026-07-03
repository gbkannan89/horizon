package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/goal/register"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"live","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready"}`)
	})

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { dsn = "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable" }

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	register.RegisterRoutes(mux, pool, time.Now)

	log.Printf("Starting Goal service on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
