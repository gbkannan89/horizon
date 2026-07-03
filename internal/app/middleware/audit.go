package middleware

import (
	"log"
	"net/http"
	"time"
)

type auditEntry struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Duration string `json:"duration"`
	UserID   string `json:"user_id,omitempty"`
	Time     string `json:"time"`
}

func AuditLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)

		entry := auditEntry{
			Method:   r.Method,
			Path:     r.URL.Path,
			Status:   sw.status,
			Duration: time.Since(start).Round(time.Millisecond).String(),
			Time:     start.UTC().Format(time.RFC3339),
		}
		if uid := r.URL.Query().Get("user_id"); uid != "" {
			entry.UserID = uid
		}

		log.Printf("AUDIT method=%s path=%s status=%d duration=%s user=%s",
			entry.Method, entry.Path, entry.Status, entry.Duration, entry.UserID)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
