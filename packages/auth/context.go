package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type contextKey string

const userKey contextKey = "user"

type User struct {
	ID          string
	Email       string
	HouseholdID string
	Roles       []string
}

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}

func GenerateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
