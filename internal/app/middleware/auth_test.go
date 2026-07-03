package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authSvc "github.com/horizon/core/services/experiences/auth/register"
)

type captureHandler struct {
	called bool
	userID string
}

func (h *captureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	if u := authSvc.GetAuthenticatedUser(r); u != nil {
		h.userID = u.UserID
	}
	w.WriteHeader(http.StatusOK)
}

func newTestJWT() *authSvc.JWTSvc {
	secret, accessTTL, refreshTTL := authSvc.DefaultJWTConfig()
	return authSvc.NewJWTSvc(secret, accessTTL, refreshTTL)
}

func TestIsPublic(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{"GET", "/health/live", true},
		{"GET", "/health/ready", true},
		{"POST", "/api/v1/auth/login", true},
		{"POST", "/api/v1/auth/refresh", true},
		{"GET", "/api/v1/protected", false},
		{"POST", "/api/v1/some/other", false},
		{"PUT", "/health/live", false},
		{"GET", "/health/", false},
		{"GET", "/health/live/extra", false},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		got := isPublic(req)
		if got != tt.want {
			t.Errorf("isPublic(%s %s) = %v, want %v", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestAuthenticate_PublicPathsPassThrough(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	for _, p := range [][2]string{
		{"GET", "/health/live"},
		{"GET", "/health/ready"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/refresh"},
	} {
		h.called = false
		rec := httptest.NewRecorder()
		mw.ServeHTTP(rec, httptest.NewRequest(p[0], p[1], nil))

		if !h.called {
			t.Errorf("public path %s %s should have passed through", p[0], p[1])
		}
		if rec.Code != http.StatusOK {
			t.Errorf("public path %s %s returned status %d, want %d", p[0], p[1], rec.Code, http.StatusOK)
		}
	}
}

func TestAuthenticate_AuthPrefixPassesThroughWithoutToken(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	for _, path := range []string{
		"/api/v1/auth/logout",
		"/api/v1/auth/me",
		"/api/v1/auth/change-password",
	} {
		h.called = false
		rec := httptest.NewRecorder()
		mw.ServeHTTP(rec, httptest.NewRequest("POST", path, nil))

		if !h.called {
			t.Errorf("auth prefix path %s should have passed through", path)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("auth prefix path %s returned status %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}

func TestAuthenticate_NoToken_Returns401(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/protected", nil))

	if h.called {
		t.Fatal("handler should NOT have been called without a token")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["success"] != false {
		t.Error("expected success=false")
	}
	errObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error object")
	}
	if errObj["code"] != "AUTHENTICATION_ERROR" {
		t.Errorf("code = %v, want AUTHENTICATION_ERROR", errObj["code"])
	}
	if !strings.Contains(errObj["message"].(string), "Authorization header") {
		t.Errorf("message = %v, want missing auth header message", errObj["message"])
	}
}

func TestAuthenticate_EmptyBearer_Returns401(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	mw.ServeHTTP(rec, req)

	if h.called {
		t.Fatal("handler should NOT have been called with empty bearer token")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthenticate_InvalidToken_Returns401(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	tokens := []string{"not-a-jwt", "a.b.c", "header.payload.wrongsig"}
	for _, tok := range tokens {
		h.called = false
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/protected", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		mw.ServeHTTP(rec, req)

		if h.called {
			t.Errorf("handler should NOT be called with invalid token %q", tok)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("token %q: status %d, want 401", tok, rec.Code)
		}
	}
}

func TestAuthenticate_ExpiredToken_Returns401(t *testing.T) {
	secret := []byte("dev-secret-change-in-production-min-32-chars!!")
	svc := authSvc.NewJWTSvc(secret, -1*time.Hour, -1*time.Hour)
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	tok, _, err := svc.GenerateTokenPair("u1", "u1@h.app", []string{"member"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	mw.ServeHTTP(rec, req)

	if h.called {
		t.Fatal("handler should NOT be called with expired token")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthenticate_RefreshTokenIsRejected(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	_, refreshTok, err := svc.GenerateTokenPair("u1", "u1@h.app", []string{"member"})
	if err != nil {
		t.Fatalf("generate tokens: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+refreshTok)
	mw.ServeHTTP(rec, req)

	if h.called {
		t.Fatal("handler should NOT be called with a refresh token")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthenticate_ValidToken_SetsContextAndPasses(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	accessTok, _, err := svc.GenerateTokenPair("user-abc-123", "test@horizon.app", []string{"member"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/protected-resource", nil)
	req.Header.Set("Authorization", "Bearer "+accessTok)
	mw.ServeHTTP(rec, req)

	if !h.called {
		t.Fatal("handler should have been called with a valid token")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if h.userID != "user-abc-123" {
		t.Errorf("userID = %q, want user-abc-123", h.userID)
	}
}

func TestAuthenticate_WithNonBearerScheme_Returns401(t *testing.T) {
	svc := newTestJWT()
	h := &captureHandler{}
	mw := Authenticate(svc)(h)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	mw.ServeHTTP(rec, req)

	if h.called {
		t.Fatal("handler should NOT be called with Basic auth")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
