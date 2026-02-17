package auth

import (
	"context"
	"net/http"

	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
)

type contextKey string

const userContextKey contextKey = "user"

var Users = map[string]*idp.User{}

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
		if !found || u.Password != idp.PasswordHash(password) {
			next.ServeHTTP(w, r)
			return
		}

		newCtx := context.WithValue(r.Context(), userContextKey, u)
		next.ServeHTTP(w, r.WithContext(newCtx))
		return
	})
}

func UserFromContext(ctx context.Context) *idp.User {
	user, ok := ctx.Value(userContextKey).(*idp.User)
	if !ok {
		return nil
	}
	return user
}
