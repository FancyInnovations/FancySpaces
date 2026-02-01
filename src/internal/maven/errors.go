package maven

import "errors"

var (
	ErrRepositoryNotFound      = errors.New("repository not found")
	ErrRepositoryAlreadyExists = errors.New("repository already exists")

	ErrArtifactNotFound      = errors.New("artifact not found")
	ErrArtifactAlreadyExists = errors.New("artifact already exists")
)
