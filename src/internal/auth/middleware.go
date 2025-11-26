package auth

import (
	"context"
	"net/http"
	"time"
)

var ApiKey string

type contextKey string

const userContextKey contextKey = "user"

var admin = User{
	ID:        "user-1",
	Provider:  ProviderBasic,
	Name:      "AdminUser",
	Email:     "admin@fancyspaces.net",
	Verified:  true,
	Password:  "...",
	Roles:     []string{"admin", "user"},
	CreatedAt: time.Now(),
	IsActive:  true,
	Metadata:  map[string]string{},
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey == ApiKey {
			newCtx := context.WithValue(r.Context(), userContextKey, &admin)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}
