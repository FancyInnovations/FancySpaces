package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

var Users = map[string]*User{}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// API key
		apiKey := r.Header.Get("Authorization")
		for _, u := range Users {
			userKey, ok := u.Metadata["api_key"]
			if !ok {
				continue
			}

			if apiKey == userKey {
				newCtx := context.WithValue(r.Context(), userContextKey, u)
				next.ServeHTTP(w, r.WithContext(newCtx))
				return
			}
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
