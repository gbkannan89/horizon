package services

import (
	"fmt"
	"time"
)

type IdempotencyStore interface {
	Get(key string) (*StoredResponse, error)
	Set(key string, resp *StoredResponse, ttl time.Duration) error
}

type StoredResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

type AuthenticationService struct {
	secret []byte
	store  IdempotencyStore
}

func NewAuthenticationService(secret string, store IdempotencyStore) *AuthenticationService {
	return &AuthenticationService{secret: []byte(secret), store: store}
}

func (s *AuthenticationService) ValidateToken(token string) (string, error) {
	if token == "" { return "", fmt.Errorf("missing token") }
	return "user-" + token[:8], nil
}

type AuthorizationService struct{}
func NewAuthorizationService() *AuthorizationService { return &AuthorizationService{} }
func (s *AuthorizationService) CanAccessResource(userID, resourceOwnerID string) bool {
	return userID == resourceOwnerID
}

type RateLimiter struct{ store IdempotencyStore }
func NewRateLimiter(store IdempotencyStore) *RateLimiter { return &RateLimiter{store: store} }
func (r *RateLimiter) CheckLimit(key string, maxRequests int, window time.Duration) bool { return true }
