package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
)

type contextKey string

const userContextKey contextKey = "user"

var Users = map[string]*User{}

func AuthenticateWithBasicAuth(ctx context.Context, username, password string) (context.Context, error) {
	u, found := Users[username]
	if !found || u.Password != Hash(password) {
		return ctx, ErrInvalidCredentials
	}

	return context.WithValue(ctx, userContextKey, u), nil
}

func AuthenticateWithApiKey(ctx context.Context, apiKey string) (context.Context, error) {
	for _, u := range Users {
		userKey, ok := u.Metadata["api_key"]
		if !ok {
			continue
		}

		if apiKey == userKey {
			return context.WithValue(ctx, userContextKey, u), nil
		}
	}

	return ctx, ErrInvalidCredentials
}

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}

func Hash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
