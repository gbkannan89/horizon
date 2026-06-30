#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 2 ]; then
  echo "Usage: $0 <type> <name>"
  echo "  type: domain, engine, infra"
  echo "  name: service-name (e.g., goal, projection)"
  exit 1
fi

TYPE=$1
NAME=$2
BASE="services/$TYPE/$NAME"

if [ -d "$BASE" ]; then
  echo "Error: Service already exists at $BASE"
  exit 1
fi

echo "Creating $TYPE service: $NAME"

mkdir -p "$BASE/cmd/server"
mkdir -p "$BASE/internal/domain"
mkdir -p "$BASE/internal/application/command"
mkdir -p "$BASE/internal/application/query"
mkdir -p "$BASE/internal/application/event"
mkdir -p "$BASE/internal/application/dto"
mkdir -p "$BASE/internal/infrastructure/persistence"
mkdir -p "$BASE/internal/infrastructure/eventbus"
mkdir -p "$BASE/internal/infrastructure/config"
mkdir -p "$BASE/internal/api/http"
mkdir -p "$BASE/spec"

# go.mod
cat > "$BASE/go.mod" <<GOMOD
module github.com/horizon/core/$BASE

go 1.24
GOMOD

# main.go
cat > "$BASE/cmd/server/main.go" <<MAIN
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"live"}`)
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready"}`)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting $NAME service on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
MAIN

echo "Service created at $BASE"
echo "Run 'cd $BASE && go build ./...' to verify compilation."
