// Package fake provides an in-memory implementation of the maven DB interface
// for tests.
package fake

import (
	"context"
	"sync"

	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
)

type DB struct {
	Mu           *sync.Mutex
	Repositories map[string]*maven.Repository
	Artifacts    map[string]*maven.Artifact
}

func New() *DB {
	return &DB{
		Mu:           &sync.Mutex{},
		Repositories: map[string]*maven.Repository{},
		Artifacts:    map[string]*maven.Artifact{},
	}
}

func (db *DB) GetRepository(ctx context.Context, spaceID, repoName string) (*maven.Repository, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName
	if r, ok := db.Repositories[key]; ok {
		cp := *r
		return &cp, nil
	}
	return nil, maven.ErrRepositoryNotFound
}

func (db *DB) GetRepositories(ctx context.Context, spaceID string) ([]maven.Repository, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var out []maven.Repository
	for _, r := range db.Repositories {
		if r.SpaceID == spaceID {
			out = append(out, *r)
		}
	}
	return out, nil
}

func (db *DB) CreateRepository(ctx context.Context, repo maven.Repository) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := repo.SpaceID + ":" + repo.Name
	if _, ok := db.Repositories[key]; ok {
		return maven.ErrRepositoryAlreadyExists
	}
	cp := repo
	db.Repositories[key] = &cp
	return nil
}

func (db *DB) UpdateRepository(ctx context.Context, repo maven.Repository) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := repo.SpaceID + ":" + repo.Name
	if _, ok := db.Repositories[key]; !ok {
		return maven.ErrRepositoryNotFound
	}
	cp := repo
	db.Repositories[key] = &cp
	return nil
}

func (db *DB) DeleteRepository(ctx context.Context, spaceID, repoName string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName
	delete(db.Repositories, key)
	return nil
}

func (db *DB) GetArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) (*maven.Artifact, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName + ":" + groupID + ":" + artifactID
	if a, ok := db.Artifacts[key]; ok {
		cp := *a
		return &cp, nil
	}
	return nil, maven.ErrArtifactNotFound
}

func (db *DB) GetArtifacts(ctx context.Context, spaceID, repoName string) ([]maven.Artifact, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var out []maven.Artifact
	for _, a := range db.Artifacts {
		if a.SpaceID == spaceID && a.Repository == repoName {
			cp := *a
			out = append(out, cp)
		}
	}
	return out, nil
}

func (db *DB) CreateArtifact(ctx context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName + ":" + artifact.Group + ":" + artifact.ID
	if _, ok := db.Artifacts[key]; ok {
		return maven.ErrArtifactAlreadyExists
	}
	cp := artifact
	db.Artifacts[key] = &cp
	return nil
}

func (db *DB) UpdateArtifact(ctx context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName + ":" + artifact.Group + ":" + artifact.ID
	if _, ok := db.Artifacts[key]; !ok {
		return maven.ErrArtifactNotFound
	}
	cp := artifact
	db.Artifacts[key] = &cp
	return nil
}

func (db *DB) DeleteArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	key := spaceID + ":" + repoName + ":" + groupID + ":" + artifactID
	delete(db.Artifacts, key)
	return nil
}
