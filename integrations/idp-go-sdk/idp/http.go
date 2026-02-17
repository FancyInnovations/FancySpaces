package idp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func (s *Service) validateRequest(r *http.Request) (*User, error) {
	if r.Header.Get("Authorization") == "" {
		return nil, ErrMissingAuthorizationHeader
	}

	authValue := r.Header.Get("Authorization")

	if strings.HasPrefix(authValue, "Bearer ") {
		token, err := tokenFromHeader(r)
		if err != nil {
			return nil, err
		}

		userID, err := s.validateToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}

		userFromCache, err := s.usersCache.GetByID(userID)
		if err == nil {
			return userFromCache, nil
		}

		resp, err := s.broker.Request("idp.user.get", []byte(userID))
		if err != nil {
			return nil, fmt.Errorf("failed to get user data: %w", err)
		}

		var u User
		if err := json.Unmarshal(resp.Data, &u); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
		}

		s.usersCache.UpsertUser(&u)

		return &u, nil
	}

	if strings.HasPrefix(authValue, "Basic ") {
		userid, password, err := basicFromHeader(r)
		if err != nil {
			return nil, err
		}

		resp, err := s.broker.Request("idp.user.validate", []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, userid, password)))
		if err != nil {
			return nil, fmt.Errorf("failed to validate user: %w", err)
		}

		var u User
		if err := json.Unmarshal(resp.Data, &u); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
		}

		return &u, nil
	}

	return nil, ErrInvalidAuthenticationMethod
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// check if the request path matches any of the excluded routes
		for _, route := range s.excludedRoutes {
			matched, err := regexp.MatchString(route, r.URL.Path)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if matched {
				next.ServeHTTP(w, r)
				return
			}
		}

		// validate the request and get the user
		user, err := s.validateRequest(r)
		if err != nil {
			fmt.Printf("Error validating user: %v\n", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
