package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

// JWTClaims represents the minimal JWT payload needed for user identification.
type JWTClaims struct {
	UserID string `json:"uid"`
}

// UserIDFromRequest extracts the authenticated user's ID from the request.
// Priority: ?user_id= query param > JWT Bearer token uid claim > "default".
func UserIDFromRequest(r *http.Request) string {
	if u := r.URL.Query().Get("user_id"); u != "" {
		return u
	}

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "default"
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "default"
	}

	// Parse JWT payload (second segment) to extract uid claim
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return token // fallback: use raw token
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return token
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return token
	}

	if claims.UserID != "" {
		return claims.UserID
	}
	return token
}
