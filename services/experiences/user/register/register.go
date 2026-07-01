package register

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/users/me", func(w http.ResponseWriter, r *http.Request) {
		handleGetProfile(w, r, pool)
	})
	mux.HandleFunc("PUT /api/v1/users/me", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateProfile(w, r, pool)
	})
	mux.HandleFunc("GET /api/v1/users/preferences", func(w http.ResponseWriter, r *http.Request) {
		handleGetPreferences(w, r, pool)
	})
	mux.HandleFunc("PUT /api/v1/users/preferences", func(w http.ResponseWriter, r *http.Request) {
		handleUpdatePreferences(w, r, pool)
	})
	mux.HandleFunc("GET /api/v1/users/privacy", func(w http.ResponseWriter, r *http.Request) {
		handleGetPrivacy(w, r, pool)
	})
	mux.HandleFunc("PUT /api/v1/users/privacy", func(w http.ResponseWriter, r *http.Request) {
		handleUpdatePrivacy(w, r, pool)
	})
}

func userID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
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

func handleGetProfile(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var id, email, name, phone, country, currency, lang, tz, createdAt string
	err := pool.QueryRow(r.Context(),
		`SELECT user_id, email, display_name, preferred_name, country, base_currency, locale, timezone, created_at FROM users WHERE user_id = $1`, uid,
	).Scan(&id, &email, &name, &phone, &country, &currency, &lang, &tz, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "User not found"); return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"id": id, "email": email, "name": name, "phone": phone,
			"country": country, "base_currency": currency, "language": lang, "timezone": tz, "created_at": createdAt,
		}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var req struct {
		Name     *string `json:"name"`
		Phone    *string `json:"phone"`
		Country  *string `json:"country"`
		Currency *string `json:"base_currency"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body"); return
	}
	_, err := pool.Exec(r.Context(),
		`UPDATE users SET display_name = COALESCE($2, display_name), preferred_name = COALESCE($3, preferred_name), country = COALESCE($4, country), base_currency = COALESCE($5, base_currency) WHERE user_id = $1`,
		uid, req.Name, req.Phone, req.Country, req.Currency)
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to update"); return }
	handleGetProfile(w, r, pool)
}

func handleGetPreferences(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	rows, err := pool.Query(r.Context(),
		`SELECT preference_key, preference_value FROM user_preferences WHERE user_id = $1`, uid)
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to load preferences"); return }
	defer rows.Close()
	prefs := map[string]string{}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		prefs[k] = v
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"preferences": prefs, "base_currency": prefs["currency"], "language": prefs["locale"],
			"theme": prefs["theme"], "email_notifications": true, "push_notifications": true,
		}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func handleUpdatePreferences(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body"); return
	}
	for k, v := range req {
		val, _ := json.Marshal(v)
		_, err := pool.Exec(r.Context(),
			`INSERT INTO user_preferences (user_id, preference_key, preference_value) VALUES ($1, $2, $3)
			 ON CONFLICT (user_id, preference_key) DO UPDATE SET preference_value = $3`,
			uid, k, string(val))
		if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to save preference"); return }
	}
	handleGetPreferences(w, r, pool)
}

func handleGetPrivacy(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var sharing, visibility string
	var thirdParty bool
	err := pool.QueryRow(r.Context(),
		`SELECT data_sharing, household_visibility, third_party_access FROM user_privacy WHERE user_id = $1`, uid,
	).Scan(&sharing, &visibility, &thirdParty)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true, "data": map[string]interface{}{
				"data_sharing_enabled": true, "analytics_enabled": true, "personalize_enabled": true,
				"third_party_sharing": false, "marketing_enabled": false,
			}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"data_sharing_enabled": sharing == "OptIn", "analytics_enabled": true,
			"personalize_enabled": true, "third_party_sharing": thirdParty, "marketing_enabled": false,
		}, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func handleUpdatePrivacy(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	uid := userID(r)
	if uid == "default" { writeError(w, http.StatusUnauthorized, "AUTH_ERROR", "Authentication required"); return }
	var req struct {
		DataSharing *bool `json:"data_sharing_enabled"`
		ThirdParty  *bool `json:"third_party_sharing"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body"); return
	}
	sharing := "OptOut"
	if req.DataSharing != nil && *req.DataSharing { sharing = "OptIn" }
	tp := false
	if req.ThirdParty != nil { tp = *req.ThirdParty }
	_, err := pool.Exec(r.Context(),
		`INSERT INTO user_privacy (user_id, data_sharing, household_visibility, third_party_access)
		 VALUES ($1, $2, 'RoleBased', $3)
		 ON CONFLICT (user_id) DO UPDATE SET data_sharing = $2, third_party_access = $3`,
		uid, sharing, tp)
	if err != nil { writeError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to update privacy"); return }
	handleGetPrivacy(w, r, pool)
}
