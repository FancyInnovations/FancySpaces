package spaces

import (
	"errors"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
)

var (
	ErrSpaceNotFound      = errors.New("space not found")
	ErrSpaceAlreadyExists = errors.New("space with given ID or slug already exists")

	ErrSlugTooLong        = errors.New("slug exceeds maximum length of 20 characters")
	ErrSlugTooShort       = errors.New("slug must be at least 3 characters long")
	ErrTitleTooLong       = errors.New("title exceeds maximum length of 100 characters")
	ErrTitleTooShort      = errors.New("title must be at least 3 characters long")
	ErrDescriptionTooLong = errors.New("description exceeds maximum length of 500 characters")

	ErrUserNotActive   = errors.New("user is not active")
	ErrUserNotVerified = errors.New("user email is not verified")
)

func ProblemFeatureNotEnabled(feature string) *problems.Problem {
	return &problems.Problem{
		Type:      "FeatureNotEnabled",
		Title:     "Feature Not Enabled",
		Detail:    "The feature '" + feature + "' is not enabled for this space.",
		Status:    http.StatusForbidden,
		Timestamp: time.Now(),
	}
}
