package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/auth"
	"github.com/horizon/core/services/infra/api-gateway/internal/errors"
)

const ContextKeyUser = "authenticated_user"

var globalAuthService *auth.Service

// InitAuth sets the global auth service reference used by the Authentication middleware.
func InitAuth(svc *auth.Service) {
	globalAuthService = svc
}

// Authentication returns middleware that validates JWT bearer tokens.
func Authentication(c *gin.Context) {
	if globalAuthService == nil {
		errors.Respond(c, 500, "INTERNAL_ERROR", "AUTH_NOT_CONFIGURED", "Authentication service not configured")
		c.Abort()
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errors.Respond(c, 401, "AUTHENTICATION_ERROR", "MISSING_TOKEN", "Authorization header is required")
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		errors.Respond(c, 401, "AUTHENTICATION_ERROR", "INVALID_TOKEN_FORMAT", "Authorization header must be: Bearer <token>")
		c.Abort()
		return
	}

	user, err := globalAuthService.Authenticate(parts[1])
	if err != nil {
		errors.Respond(c, 401, "AUTHENTICATION_ERROR", "TOKEN_INVALID", err.Error())
		c.Abort()
		return
	}

	c.Set(ContextKeyUser, user)
	c.Next()
}

// Authorization returns middleware that checks if the authenticated user has the required role.
func Authorization(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRaw, exists := c.Get(ContextKeyUser)
		if !exists {
			errors.Respond(c, 401, "AUTHENTICATION_ERROR", "NOT_AUTHENTICATED", "Authentication required")
			c.Abort()
			return
		}

		user, ok := userRaw.(*auth.AuthenticatedUser)
		if !ok {
			errors.Respond(c, 500, "INTERNAL_ERROR", "INVALID_CONTEXT", "Invalid authentication context")
			c.Abort()
			return
		}

		if len(requiredRoles) == 0 {
			c.Next()
			return
		}

		roleMap := make(map[string]bool)
		for _, role := range user.Roles {
			roleMap[role] = true
		}

		if roleMap["admin"] {
			c.Next()
			return
		}

		for _, required := range requiredRoles {
			if roleMap[required] {
				c.Next()
				return
			}
		}

		errors.Respond(c, 403, "AUTHORIZATION_ERROR", "INSUFFICIENT_PERMISSIONS", "You do not have permission to perform this action")
		c.Abort()
	}
}

// CurrentUser extracts the authenticated user from the context.
func CurrentUser(c *gin.Context) *auth.AuthenticatedUser {
	userRaw, exists := c.Get(ContextKeyUser)
	if !exists {
		return nil
	}
	user, ok := userRaw.(*auth.AuthenticatedUser)
	if !ok {
		return nil
	}
	return user
}
