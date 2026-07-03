package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/household/internal/application"
	"github.com/horizon/core/services/domains/household/internal/domain"
	"github.com/horizon/core/services/domains/household/internal/infrastructure/persistence"
)

type noopPublisher struct{}

func (n *noopPublisher) Publish(event domain.DomainEvent) error { return nil }

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8110" }

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { dsn = "postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable" }

	pool, err := pgxpool.New(nil, dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	repo := persistence.NewHouseholdRepository(pool)
	publisher := &noopPublisher{}
	svc := application.NewHouseholdService(repo, publisher, time.Now)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"live"}`))
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(503)
			w.Write([]byte(`{"status":"not ready"}`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ready"}`))
	})

	_ = svc // service wired, handlers added in future milestones

	log.Printf("Household domain service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
