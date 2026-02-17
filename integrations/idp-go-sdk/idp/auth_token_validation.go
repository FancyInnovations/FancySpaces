package idp

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const ServiceName = "fancyanalytics-idp"

var SigningMethod = jwt.SigningMethodRS256

// validateToken validates the given JWT token string and returns the user ID if the token is valid.
func (s *Service) validateToken(tokenString string) (string, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{SigningMethod.Alg()}),
		jwt.WithIssuer(ServiceName),
	)

	token, err := parser.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		s.tokenKeyFunc,
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	if claims.Issuer != ServiceName {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}

// tokenKeyFunc is a helper function that returns the public key for validating the token's signature.
func (s *Service) tokenKeyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != SigningMethod.Alg() {
		return nil, fmt.Errorf("unexpected signing method")
	}

	return s.publicKey, nil
}
