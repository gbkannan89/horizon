package register

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Handlers struct {
	authService *JWTSvc
}

func New() *Handlers {
	secret, accessTTL, refreshTTL := DefaultJWTConfig()
	return &Handlers{authService: NewJWTSvc(secret, accessTTL, refreshTTL)}
}

func NewWithSvc(svc *JWTSvc) *Handlers {
	return &Handlers{authService: svc}
}

func RegisterRoutes(mux *http.ServeMux) {
	h := New()
	RegisterRoutesWithSvc(mux, h.authService)
}

func RegisterRoutesWithSvc(mux *http.ServeMux, svc *JWTSvc) {
	h := NewWithSvc(svc)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", h.Logout)
	mux.Handle("GET /api/v1/auth/me", Middleware(h.authService)(http.HandlerFunc(h.Me)))
	mux.Handle("POST /api/v1/auth/change-password", Middleware(h.authService)(http.HandlerFunc(h.ChangePassword)))
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required")
		return
	}
	userID := req.Email
	if req.Email == "demo@horizon.app" {
		userID = "a1b2c3d4-0000-4000-8000-000000000001"
	}
	roles := []string{"member"}
	accessToken, refreshToken, err := h.authService.GenerateTokenPair(userID, req.Email, roles)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate tokens")
		return
	}
	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: accessToken, RefreshToken: refreshToken,
		ExpiresIn: int(15 * time.Minute / time.Second), TokenType: "Bearer",
	})
}

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct{ RefreshToken string `json:"refresh_token"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body"); return
	}
	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "refresh_token is required"); return
	}
	newAccess, newRefresh, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Invalid or expired refresh token"); return
	}
	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: newAccess, RefreshToken: newRefresh,
		ExpiresIn: int(15 * time.Minute / time.Second), TokenType: "Bearer",
	})
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// ChangePassword updates the user's password. MVP: accepts any valid request.
func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := GetAuthenticatedUser(r)
	if user == nil {
		writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Authentication required")
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Both current and new password are required")
		return
	}
	if len(req.NewPassword) < 8 {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "New password must be at least 8 characters")
		return
	}
	// MVP: accept any valid password. Phase 3+ will implement real password hashing.
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Password changed successfully"})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	user := GetAuthenticatedUser(r)
	if user == nil {
		writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Authentication required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": user.UserID, "email": user.Email, "roles": user.Roles,
	})
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

type contextKey string

const userKey contextKey = "authenticated_user"

func Middleware(authService *JWTSvc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Missing or invalid Authorization header")
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := authService.Authenticate(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Invalid or expired token")
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthenticatedUser(r *http.Request) *AuthenticatedUser {
	user, _ := r.Context().Value(userKey).(*AuthenticatedUser)
	return user
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Correlation-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Correlation-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w, r)
	})
}
