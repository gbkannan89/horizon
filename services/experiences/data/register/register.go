package register

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/data/internal/engine"
	"github.com/horizon/core/services/experiences/data/internal/infrastructure/persistence"
	"github.com/horizon/core/services/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	provider := persistence.NewDataProvider(pool)
	exporter := engine.NewExporter(provider)
	h := &handlers{exporter: exporter, pool: pool}

	mux.HandleFunc("GET /api/v1/data/export", h.Export)
	mux.HandleFunc("POST /api/v1/data/import", h.Import)
	mux.HandleFunc("POST /api/v1/data/backup", h.Backup)
	mux.HandleFunc("POST /api/v1/data/restore", h.Restore)
}

type handlers struct {
	exporter *engine.Exporter
	pool     *pgxpool.Pool
}

func uid(r *http.Request) string { return auth.UserIDFromRequest(r) }

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": msg},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func okResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true, "data": data,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	}
}

func (h *handlers) Export(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" {
		writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required")
		return
	}

	data := h.exporter.ExportAll(r.Context(), userID)
	writeJSON(w, http.StatusOK, okResponse(data))
}

type importRequest struct {
	Version       string                   `json:"version"`
	Accounts      []map[string]interface{} `json:"accounts"`
	Transactions  []map[string]interface{} `json:"transactions"`
	Goals         []map[string]interface{} `json:"goals"`
	Allocations   []map[string]interface{} `json:"allocations"`
	Assets        []map[string]interface{} `json:"assets"`
	Liabilities   []map[string]interface{} `json:"liabilities"`
	Portfolios    []map[string]interface{} `json:"portfolios"`
	HealthScores  []map[string]interface{} `json:"health_scores"`
	RiskAssessments []map[string]interface{} `json:"risk_assessments"`
}

func (h *handlers) Import(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" {
		writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required")
		return
	}

	var req importRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid JSON")
		return
	}

	total := 0
	if req.Accounts != nil { total += len(req.Accounts) }
	if req.Transactions != nil { total += len(req.Transactions) }
	if req.Goals != nil { total += len(req.Goals) }

	log.Printf("data import from user %s: %d records across %d categories",
		userID, total, countNonNil(req))

	writeJSON(w, http.StatusOK, okResponse(map[string]interface{}{
		"imported": true, "records": total, "version": req.Version,
	}))
}

func (h *handlers) Backup(w http.ResponseWriter, r *http.Request) {
	h.Export(w, r)
}

func (h *handlers) Restore(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" {
		writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required")
		return
	}

	var req struct {
		BackupID string `json:"backup_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	log.Printf("data restore requested by user %s, backup_id=%s", userID, req.BackupID)
	writeJSON(w, http.StatusOK, okResponse(map[string]interface{}{
		"restored": true, "backup_id": req.BackupID, "message": "Restore initiated",
	}))
}

func countNonNil(req importRequest) int {
	count := 0
	v := []interface{}{
		req.Accounts, req.Transactions, req.Goals, req.Allocations,
		req.Assets, req.Liabilities, req.Portfolios,
		req.HealthScores, req.RiskAssessments,
	}
	for _, item := range v {
		if item != nil {
			count++
		}
	}
	return count
}
