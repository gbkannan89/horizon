package main

import (
	"log"
	"net/http"
	"os"

	"github.com/horizon/core/services/engines/tax-planning/register"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8160" }

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"live"}`))
	})
	register.RegisterRoutes(mux)
	log.Printf("Tax planning service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
