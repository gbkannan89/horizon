package api

import (
	"log"
	"net/http"
	"time"
)

func NewRouter(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()
	handlers.Register(mux)
	mux.HandleFunc("GET /health/live", healthLive)
	mux.HandleFunc("GET /health/ready", healthReady)
	return recovery(requestLog(mux))
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"success": false, "error": map[string]string{"code": "PANIC", "message": "internal error"},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func healthLive(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "live"})
}

func healthReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
