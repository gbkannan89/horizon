package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"
const UserIDKey contextKey = "user_id"

// CorrelationID generates or forwards a correlation ID.
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Header.Get("X-Correlation-ID")
		if cid == "" {
			b := make([]byte, 16)
			rand.Read(b)
			cid = hex.EncodeToString(b)
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, cid)
		w.Header().Set("X-Correlation-ID", cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestLog logs every request with duration.
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Recovery recovers from panics.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v", rec)
				http.Error(w, `{"error":"internal_server_error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
