package spaces

import (
	"errors"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
)

var (
	ErrSpaceAlreadyExists = errors.New("space with given ID or slug already exists")
	ErrUserNotActive      = errors.New("user is not active")
	ErrUserNotVerified    = errors.New("user email is not verified")
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
