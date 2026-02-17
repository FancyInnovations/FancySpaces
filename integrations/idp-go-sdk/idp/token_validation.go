package idp

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const ServiceName = "fancyanalytics-idp"

var SigningMethod = jwt.SigningMethodRS256

// ValidateToken validates the provided JWT token string and returns the associated user if the token is valid.
func (s *Service) ValidateToken(token string) (*User, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{SigningMethod.Alg()}),
		jwt.WithIssuer(ServiceName),
	)

	t, err := parser.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		s.tokenKeyFunc,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok || !t.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Issuer != ServiceName {
		return nil, ErrInvalidToken
	}

	userID := claims.Subject

	return s.GetUser(userID)
}

// tokenKeyFunc is a helper function that returns the public key for validating the token's signature.
func (s *Service) tokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != SigningMethod.Alg() {
		return nil, fmt.Errorf("unexpected signing method")
	}

	return s.publicKey, nil
}
