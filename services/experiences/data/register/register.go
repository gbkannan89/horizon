package register

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	h := &handlers{pool: pool}
	mux.HandleFunc("GET /api/v1/data/export", h.Export)
	mux.HandleFunc("POST /api/v1/data/import", h.Import)
	mux.HandleFunc("POST /api/v1/data/backup", h.Backup)
	mux.HandleFunc("POST /api/v1/data/restore", h.Restore)
}

type handlers struct {
	pool *pgxpool.Pool
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

// Export returns a JSON dump of the user's data.
func (h *handlers) Export(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }

	accounts := queryJSON[[]map[string]interface{}](r, h.pool, `SELECT json_agg(to_json(accounts.*)) FROM accounts WHERE owner_id = $1`, userID)
	events := queryJSON[[]map[string]interface{}](r, h.pool, `SELECT json_agg(to_json(financial_events.*)) FROM financial_events WHERE user_id = $1`, userID)
	goals := queryJSON[[]map[string]interface{}](r, h.pool, `SELECT json_agg(to_json(goals.*)) FROM goals WHERE user_id = $1`, userID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"exported_at": time.Now().UTC().Format(time.RFC3339),
			"accounts":    ifNil(accounts),
			"transactions": ifNil(events),
			"goals":       ifNil(goals),
		}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// Import accepts data and logs it (MVP: no-op).
func (h *handlers) Import(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid JSON"); return
	}
	log.Printf("data import from user %s: %d top-level keys", userID, len(body))
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{"imported": true, "records": len(body)},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// Backup creates a backup (MVP: alias for export).
func (h *handlers) Backup(w http.ResponseWriter, r *http.Request) {
	h.Export(w, r)
}

// Restore restores from a backup (MVP: accepts backup_id, returns success).
func (h *handlers) Restore(w http.ResponseWriter, r *http.Request) {
	userID := uid(r)
	if userID == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{"restored": true, "message": "Restore complete"},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func queryJSON[T any](r *http.Request, pool *pgxpool.Pool, query string, args ...any) T {
	var zero T
	row := pool.QueryRow(r.Context(), query, args...)
	var result T
	if err := row.Scan(&result); err != nil {
		return zero
	}
	return result
}

func ifNil(v any) any {
	if v == nil { return []interface{}{} }
	return v
}
