package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server wraps an HTTP server with graceful shutdown.
type Server struct {
	server *http.Server
}

// NewServer creates an HTTP server with sensible timeouts.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// StartAndWait starts the server and blocks until SIGINT/SIGTERM.
func (s *Server) StartAndWait() {
	go func() {
		log.Printf("server listening on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
}

// RecoveryMiddleware recovers from panics and returns a JSON error response.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"success":false,"error":{"code":"PANIC","message":"internal error"}}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// HealthHandler returns a simple health check handler.
func HealthHandler(status string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"` + status + `"}`))
	}
}
