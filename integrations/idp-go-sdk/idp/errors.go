package idp

import (
	"errors"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
)

var (
	ErrMissingAuthorizationHeader  = errors.New("missing Authorization header")
	ErrInvalidAuthenticationMethod = errors.New("invalid authentication method, expected Bearer or Basic")

	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidToken       = errors.New("invalid token")

	ErrInvalidApiKeyFormat = errors.New("invalid API key format")
	ErrApiKeyNotFound      = errors.New("API key not found")
	ErrInvalidApiKey       = errors.New("invalid API key")

	ErrInvalidBasicCredentials = errors.New("invalid basic authentication credentials")

	ErrUserNotFound = errors.New("user not found")
)

func AccountNotVerifiedProblem() *problems.Problem {
	return &problems.Problem{
		Type:      "AccountNotVerified",
		Title:     "Account Not Verified",
		Detail:    "Your account is not verified. Please verify your account to access this feature.",
		Status:    http.StatusForbidden,
		Timestamp: time.Now(),
	}
}

func AccountDisabledProblem() *problems.Problem {
	return &problems.Problem{
		Type:      "AccountDisabled",
		Title:     "Account Disabled",
		Detail:    "Your account has been disabled. Please contact support for assistance.",
		Status:    http.StatusForbidden,
		Timestamp: time.Now(),
	}
}
