package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// Check represents a single health check result.
type Check struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // healthy, degraded, unhealthy
	Message string `json:"message,omitempty"`
	Latency string `json:"latency_ms,omitempty"`
}

// Handler provides HTTP health check endpoints.
type Handler struct {
	checks []func() Check
}

// NewHandler creates a new health check handler.
func NewHandler() *Handler { return &Handler{} }

// Register adds a health check function.
func (h *Handler) Register(fn func() Check) { h.checks = append(h.checks, fn) }

// Live returns the liveness endpoint handler.
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"status": "live", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// Ready returns the readiness endpoint handler.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	overall := "healthy"
	results := make([]Check, 0, len(h.checks))
	for _, fn := range h.checks {
		c := fn()
		results = append(results, c)
		if c.Status != "healthy" {
			overall = "degraded"
		}
	}
	status := http.StatusOK
	if overall != "healthy" {
		status = http.StatusServiceUnavailable
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": overall,
		"checks": results,
	})
}
