package main

import (
	"log"
	"net/http"
	"os"

	"github.com/horizon/core/services/engines/documents/register"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8180" }
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"live"}`))
	})
	register.RegisterRoutes(mux)
	log.Printf("Documents service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
