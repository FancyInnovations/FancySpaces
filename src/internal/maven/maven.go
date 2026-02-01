package maven

import (
	"context"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/internal/analytics"
)

type DB interface {
}

type Store struct {
	db        DB
	analytics *analytics.Store
}

type Configuration struct {
	DB        DB
	Analytics *analytics.Store
}

func New(cfg Configuration) *Store {
	return &Store{
		db:        cfg.DB,
		analytics: cfg.Analytics,
	}
}

func (s *Store) GetRepository(ctx context.Context, spaceID, repoName string) (*Repository, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Store) GetRepositories(ctx context.Context, spaceID string) ([]Repository, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Store) CreateRepository(ctx context.Context, repo Repository) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) UpdateRepository(ctx context.Context, repo Repository) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) DeleteRepository(ctx context.Context, spaceID, repoName string) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) GetArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) (*Artifact, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Store) GetArtifacts(ctx context.Context, spaceID, repoName string) ([]Artifact, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Store) CreateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) UpdateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) DeleteArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) UploadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string, data []byte) error {
	return fmt.Errorf("not implemented")
}

func (s *Store) DownloadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
