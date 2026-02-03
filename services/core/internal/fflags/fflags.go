package fflags

import "github.com/OliverSchlueter/goutils/featureflags"

var (
	DisableIssueSyncer = featureflags.Register("DISABLE_ISSUE_SYNCER_FEATURE_FLAG")
)
