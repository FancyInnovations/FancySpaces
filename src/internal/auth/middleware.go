package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

var ApiKey string

type contextKey string

const userContextKey contextKey = "user"

var UserAdmin = User{
	ID:        "admin-1",
	Provider:  ProviderBasic,
	Name:      "Admin",
	Email:     "admin@fancyspaces.net",
	Verified:  true,
	Password:  "something",
	Roles:     []string{"admin", "user"},
	CreatedAt: time.Date(2025, 12, 3, 19, 0, 0, 0, time.UTC),
	IsActive:  true,
	Metadata:  map[string]string{},
}

var Users = map[string]*User{
	UserAdmin.ID: &UserAdmin,
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// API key
		apiKey := r.Header.Get("Authorization")
		if apiKey == ApiKey {
			newCtx := context.WithValue(r.Context(), userContextKey, &UserAdmin)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		// Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		u, found := Users[username]
		if !found || u.Password != Hash(password) {
			next.ServeHTTP(w, r)
			return
		}

		newCtx := context.WithValue(r.Context(), userContextKey, u)
		next.ServeHTTP(w, r.WithContext(newCtx))
		return
	})
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
