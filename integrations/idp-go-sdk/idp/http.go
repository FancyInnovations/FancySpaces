package idp

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/OliverSchlueter/goutils/sloki"
)

func (s *Service) validateHTTPRequest(r *http.Request) (*User, error) {
	if r.Header.Get("Authorization") == "" {
		return nil, ErrMissingAuthorizationHeader
	}

	authValue := r.Header.Get("Authorization")

	if strings.HasPrefix(authValue, "Bearer ") {
		token, err := tokenFromHeader(r)
		if err != nil {
			return nil, err
		}

		u, err := s.ValidateToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}
		return u, nil
	}

	if strings.HasPrefix(authValue, "Basic ") {
		userid, password, err := basicFromHeader(r)
		if err != nil {
			return nil, err
		}

		u, err := s.ValidateUser(userid, password)
		if err != nil {
			return nil, fmt.Errorf("failed to validate basic credentials: %w", err)
		}

		return u, nil
	}

	return nil, ErrInvalidAuthenticationMethod
}

func (s *Service) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// validate the request and get the user
		user, err := s.validateHTTPRequest(r)
		if err != nil {
			slog.Warn("Unauthorized request", sloki.WrapError(err))
			next.ServeHTTP(w, r)
			return
		}

		// attach the user to the context
		ctx := attachUserToCtx(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// tokenFromHeader extracts the token from the Authorization header and validates its format.
func tokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidTokenFormat
	}

	return authHeader[len("Bearer "):], nil
}

// basicFromHeader extracts the username and password from the Authorization header and validates its format.
func basicFromHeader(r *http.Request) (userid, password string, err error) {
	if r.Header.Get("Authorization") == "" {
		return "", "", ErrMissingAuthorizationHeader
	}

	userid, password, ok := r.BasicAuth()
	if !ok {
		return "", "", ErrInvalidTokenFormat
	}

	if userid == "" || password == "" {
		return "", "", ErrInvalidTokenFormat
	}

	return userid, password, nil
}
