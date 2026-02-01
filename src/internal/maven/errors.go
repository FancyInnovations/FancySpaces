package maven

import "errors"

var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrArtifactNotFound   = errors.New("artifact not found")
)
