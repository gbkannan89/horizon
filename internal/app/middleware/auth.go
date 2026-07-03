package middleware

import (
	"net/http"
	"strings"

	authReg "github.com/horizon/core/services/experiences/auth/register"
)

var publicPaths = []string{
	"GET /health/live",
	"GET /health/ready",
	"POST /api/v1/auth/login",
	"POST /api/v1/auth/refresh",
}

func isPublic(r *http.Request) bool {
	key := r.Method + " " + r.URL.Path
	for _, p := range publicPaths {
		if key == p {
			return true
		}
	}
	return false
}

func Authenticate(authService *authReg.JWTSvc) func(http.Handler) http.Handler {
	authMw := authReg.Middleware(authService)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublic(r) {
				next.ServeHTTP(w, r)
				return
			}
			if strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
				next.ServeHTTP(w, r)
				return
			}
			authMw(next).ServeHTTP(w, r)
		})
	}
}
