package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	p := os.Getenv("PORT"); if p == "" { p = "8095" }
	m := http.NewServeMux()
	m.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200); fmt.Fprintf(w, `{"status":"live","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	})
	m.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200); fmt.Fprintf(w, `{"status":"ready"}`)
	})
	log.Printf("Starting Optimization Engine on :%s", p)
	log.Fatal(http.ListenAndServe(":"+p, m))
}
