package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

// JWTClaims represents the minimal JWT payload needed for user identification.
type JWTClaims struct {
	UserID         string            `json:"uid"`
	HouseholdRoles map[string]string `json:"household_roles,omitempty"`
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

// HouseholdRoleFromRequest extracts the user's role in a specific household.
// Returns an empty string if the user is not a member of the household.
func HouseholdRoleFromRequest(r *http.Request, householdID string) string {
	// For testing and local dev, if a role is passed in the header, use it
	if role := r.Header.Get("X-Household-Role-" + householdID); role != "" {
		return role
	}

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return ""
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}

	if claims.HouseholdRoles != nil {
		return claims.HouseholdRoles[householdID]
	}
	return ""
}
