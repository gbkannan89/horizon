package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/auth"
	"github.com/horizon/core/services/infra/api-gateway/internal/errors"
	"github.com/horizon/core/services/infra/api-gateway/internal/middleware"
)

// LoginRequest represents a login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents a successful login response.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshRequest represents a refresh token request.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login authenticates a user and returns tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.Validation(c, "Invalid request: "+err.Error())
		return
	}

	// For MVP, accept any valid email/password combination
	// Phase 5 will wire this to the User application service
	userID := "user-" + req.Email
	roles := []string{"member"}

	accessToken, refreshToken, err := h.authService.GenerateTokenPair(userID, req.Email, roles)
	if err != nil {
		errors.Internal(c, "Failed to generate tokens")
		return
	}

	respondJSON(c, http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(15 * time.Minute / time.Second),
		TokenType:    "Bearer",
	})
}

// Refresh issues a new token pair from a valid refresh token.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.Validation(c, "Invalid request: "+err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		errors.Respond(c, 401, "AUTHENTICATION_ERROR", "REFRESH_FAILED", err.Error())
		return
	}

	respondJSON(c, http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(15 * time.Minute / time.Second),
		TokenType:    "Bearer",
	})
}

// Logout invalidates the current session.
func (h *AuthHandler) Logout(c *gin.Context) {
	// For MVP, logout is client-side (discard tokens)
	// Phase 5 will implement token blacklisting via Redis
	respondNoContent(c)
}

// Me returns the current authenticated user's profile.
func (h *AuthHandler) Me(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		errors.Respond(c, 401, "AUTHENTICATION_ERROR", "NOT_AUTHENTICATED", "Authentication required")
		return
	}
	respondJSON(c, http.StatusOK, gin.H{
		"user_id": user.UserID,
		"email":   user.Email,
		"roles":   user.Roles,
	})
}
