package register

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims represents the JWT payload.
type Claims struct {
	UserID    string   `json:"uid"`
	Email     string   `json:"email,omitempty"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"ttype"`
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
}

// AuthenticatedUser represents a verified user in the request context.
type AuthenticatedUser struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

// JWTSvc handles JWT token generation and validation.
type JWTSvc struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// DefaultJWTConfig returns a sensible development configuration.
func DefaultJWTConfig() ([]byte, time.Duration, time.Duration) {
	return []byte("dev-secret-change-in-production-min-32-chars!!"),
		15 * time.Minute,
		7 * 24 * time.Hour
}

// NewJWTSvc creates a new JWT service.
func NewJWTSvc(secret []byte, accessTTL, refreshTTL time.Duration) *JWTSvc {
	return &JWTSvc{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

// GenerateTokenPair creates an access token and refresh token.
func (s *JWTSvc) GenerateTokenPair(userID, email string, roles []string) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateToken(userID, email, roles, "access", s.accessTTL)
	if err != nil { return "", "", fmt.Errorf("generate access: %w", err) }
	refreshToken, err = s.generateToken(userID, email, roles, "refresh", s.refreshTTL)
	if err != nil { return "", "", fmt.Errorf("generate refresh: %w", err) }
	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT and returns the claims.
func (s *JWTSvc) ValidateToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 { return nil, fmt.Errorf("invalid token format") }
	expectedSig := s.sign(parts[0] + "." + parts[1])
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) { return nil, fmt.Errorf("invalid signature") }
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil { return nil, fmt.Errorf("decode payload: %w", err) }
	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt+30 { return nil, fmt.Errorf("token expired") }
	return &claims, nil
}

// Authenticate validates a bearer token and returns the authenticated user.
func (s *JWTSvc) Authenticate(token string) (*AuthenticatedUser, error) {
	claims, err := s.ValidateToken(token)
	if err != nil { return nil, fmt.Errorf("authenticate: %w", err) }
	if claims.TokenType != "access" { return nil, fmt.Errorf("invalid token type: %s", claims.TokenType) }
	return &AuthenticatedUser{UserID: claims.UserID, Email: claims.Email, Roles: claims.Roles}, nil
}

// Refresh validates a refresh token and issues a new token pair.
func (s *JWTSvc) Refresh(refreshToken string) (newAccess, newRefresh string, err error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil { return "", "", fmt.Errorf("invalid refresh: %w", err) }
	if claims.TokenType != "refresh" { return "", "", fmt.Errorf("invalid token type: %s", claims.TokenType) }
	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Roles)
}

func (s *JWTSvc) generateToken(userID, email string, roles []string, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID, Email: email, Roles: roles, TokenType: tokenType,
		ExpiresAt: now.Add(ttl).Unix(), IssuedAt: now.Unix(),
	}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil { return "", err }
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payload + "." + s.sign(header+"."+payload), nil
}

func (s *JWTSvc) sign(data string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
