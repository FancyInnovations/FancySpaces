package auth

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// API key
		apiKey := r.Header.Get("Authorization")
		newCtx, err := AuthenticateWithApiKey(r.Context(), apiKey)
		if err == nil {
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		// Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		newCtx, err = AuthenticateWithBasicAuth(r.Context(), username, password)
		next.ServeHTTP(w, r.WithContext(newCtx))
		return
	})
}
