package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/internal/auth"
)

type Handlers struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Handlers {
	return &Handlers{pool: pool}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// GetProfile returns the authenticated user's profile.
func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	if userID == "default" {
		writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Authentication required")
		return
	}
	var id, email, name, phone, country, currency, lang, tz, createdAt string
	err := h.pool.QueryRow(r.Context(),
		`SELECT user_id, email, display_name, preferred_name, country, base_currency, locale, timezone, created_at
		 FROM users WHERE user_id = $1`, userID,
	).Scan(&id, &email, &name, &phone, &country, &currency, &lang, &tz, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id": id, "email": email, "name": name, "phone": phone,
			"country": country, "base_currency": currency,
			"language": lang, "timezone": tz, "created_at": createdAt,
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// UpdateProfile updates the authenticated user's profile fields.
func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromRequest(r)
	if userID == "default" {
		writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Authentication required")
		return
	}
	var req struct {
		Name    *string `json:"name"`
		Phone   *string `json:"phone"`
		Country *string `json:"country"`
		Currency *string `json:"base_currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	// Build dynamic UPDATE — for MVP, update known fields
	_, err := h.pool.Exec(r.Context(),
		`UPDATE users SET
			display_name = COALESCE($2, display_name),
			preferred_name = COALESCE($3, preferred_name),
			country = COALESCE($4, country),
			base_currency = COALESCE($5, base_currency)
		 WHERE user_id = $1`,
		userID, req.Name, req.Phone, req.Country, req.Currency,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to update profile")
		return
	}
	// Re-fetch and return updated profile
	h.GetProfile(w, r)
}
