package maven

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/analytics"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/javadoccache"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
)

type DB interface {
	GetRepository(ctx context.Context, spaceID, repoName string) (*Repository, error)
	GetRepositories(ctx context.Context, spaceID string) ([]Repository, error)
	CreateRepository(ctx context.Context, repo Repository) error
	UpdateRepository(ctx context.Context, repo Repository) error
	DeleteRepository(ctx context.Context, spaceID, repoName string) error

	GetArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) (*Artifact, error)
	GetArtifacts(ctx context.Context, spaceID, repoName string) ([]Artifact, error)
	CreateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error
	UpdateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error
	DeleteArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) error
}

type FileStorage interface {
	UploadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string, data []byte) error
	DownloadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) ([]byte, error)
	DeleteArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) error
}

type JavadocCache interface {
	CacheJavadoc(key string, javadocZipData []byte) (map[string][]byte, error)
	GetJavadocFile(key, filePath string) ([]byte, error)
	IsJavadocCached(key string) bool
}

type Store struct {
	spaces       *spaces.Store
	db           DB
	fileStore    FileStorage
	fileCache    FileStorage
	javadocCache JavadocCache
	analytics    *analytics.Store
}

type Configuration struct {
	Spaces       *spaces.Store
	DB           DB
	FileStore    FileStorage
	FileCache    FileStorage
	JavadocCache JavadocCache
	Analytics    *analytics.Store
}

func New(cfg Configuration) *Store {
	return &Store{
		spaces:       cfg.Spaces,
		db:           cfg.DB,
		fileStore:    cfg.FileStore,
		fileCache:    cfg.FileCache,
		javadocCache: cfg.JavadocCache,
		analytics:    cfg.Analytics,
	}
}

func (s *Store) GetRepository(ctx context.Context, spaceID, repoName string) (*Repository, error) {
	return s.db.GetRepository(ctx, spaceID, repoName)
}

func (s *Store) GetRepositories(ctx context.Context, spaceID string) ([]Repository, error) {
	return s.db.GetRepositories(ctx, spaceID)
}

func (s *Store) CreateRepository(ctx context.Context, spaceID string, repo Repository) error {
	_, err := s.GetRepository(ctx, spaceID, repo.Name)
	if err == nil {
		return ErrRepositoryAlreadyExists
	}

	repo.CreatedAt = time.Now()
	repo.SpaceID = spaceID

	return s.db.CreateRepository(ctx, repo)
}

func (s *Store) UpdateRepository(ctx context.Context, spaceID string, repo Repository) error {
	repo.SpaceID = spaceID

	return s.db.UpdateRepository(ctx, repo)
}

func (s *Store) DeleteRepository(ctx context.Context, spaceID, repoName string) error {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return err
	}

	if repo.InternalMirror == nil {
		artifacts, err := s.GetArtifacts(ctx, spaceID, repo.Name)
		if err != nil {
			return err
		}

		// Delete all artifacts
		for _, artifact := range artifacts {
			if err := s.DeleteArtifact(ctx, spaceID, repo.Name, artifact.Group, artifact.ID); err != nil {
				return err
			}
		}
	}

	return s.db.DeleteRepository(ctx, spaceID, repoName)
}

func (s *Store) GetArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) (*Artifact, error) {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return nil, err
	}

	// If the repository is an internal mirror, redirect the request to the mirrored repository
	if repo.InternalMirror != nil {
		return s.GetArtifact(ctx, repo.InternalMirror.SpaceID, repo.InternalMirror.Repository, groupID, artifactID)
	}

	return s.db.GetArtifact(ctx, spaceID, repoName, groupID, artifactID)
}

func (s *Store) GetArtifacts(ctx context.Context, spaceID, repoName string) ([]Artifact, error) {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return nil, err
	}

	// If the repository is an internal mirror, redirect the request to the mirrored repository
	if repo.InternalMirror != nil {
		return s.GetArtifacts(ctx, repo.InternalMirror.SpaceID, repo.InternalMirror.Repository)
	}

	return s.db.GetArtifacts(ctx, spaceID, repoName)
}

func (s *Store) CreateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return err
	}

	// If the repository is an internal mirror, forbid artifact creation
	if repo.InternalMirror != nil {
		return ErrCannotModifyArtifactInMirrorRepository
	}

	_, err = s.GetArtifact(ctx, spaceID, repoName, artifact.Group, artifact.ID)
	if err == nil {
		return ErrArtifactAlreadyExists
	}

	return s.db.CreateArtifact(ctx, spaceID, repoName, artifact)
}

func (s *Store) UpdateArtifact(ctx context.Context, spaceID, repoName string, artifact Artifact) error {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return err
	}

	// If the repository is an internal mirror, forbid artifact update
	if repo.InternalMirror != nil {
		return ErrCannotModifyArtifactInMirrorRepository
	}

	return s.db.UpdateArtifact(ctx, spaceID, repoName, artifact)
}

func (s *Store) DeleteArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) error {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return err
	}

	// If the repository is an internal mirror, forbid artifact deletion
	if repo.InternalMirror != nil {
		return ErrCannotModifyArtifactInMirrorRepository
	}

	artifact, err := s.GetArtifact(ctx, spaceID, repoName, groupID, artifactID)
	if err != nil {
		return err
	}

	// Delete all artifact files from file storage
	for _, version := range artifact.Versions {
		for _, file := range version.Files {
			if err := s.fileStore.DeleteArtifactFile(ctx, spaceID, repoName, groupID, artifactID, version.Version, file.Name); err != nil {
				return err
			}
			if err := s.fileCache.DeleteArtifactFile(ctx, spaceID, repoName, groupID, artifactID, version.Version, file.Name); err != nil {
				return err
			}
		}
	}

	return s.db.DeleteArtifact(ctx, spaceID, repoName, groupID, artifactID)
}

func (s *Store) UploadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string, data []byte) error {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return err
	}

	// If the repository is an internal mirror, forbid artifact file upload
	if repo.InternalMirror != nil {
		return ErrCannotModifyArtifactInMirrorRepository
	}

	groupPath := strings.ReplaceAll(groupID, ".", "/")

	if err := s.fileCache.UploadArtifactFile(ctx, spaceID, repoName, groupPath, artifactID, version, fileName, data); err != nil {
		return err
	}

	return s.fileStore.UploadArtifactFile(ctx, spaceID, repoName, groupPath, artifactID, version, fileName, data)
}

func (s *Store) DownloadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) ([]byte, error) {
	repo, err := s.GetRepository(ctx, spaceID, repoName)
	if err != nil {
		return nil, err
	}

	// If the repository is an internal mirror, redirect the request to the mirrored repository
	if repo.InternalMirror != nil {
		return s.DownloadArtifactFile(ctx, repo.InternalMirror.SpaceID, repo.InternalMirror.Repository, groupID, artifactID, version, fileName)
	}

	groupPath := strings.ReplaceAll(groupID, ".", "/")

	data, err := s.fileCache.DownloadArtifactFile(ctx, spaceID, repoName, groupPath, artifactID, version, fileName)
	if err == nil {
		return data, nil
	}

	return s.fileStore.DownloadArtifactFile(ctx, spaceID, repoName, groupPath, artifactID, version, fileName)
}

func (s *Store) GetJavadocFile(ctx context.Context, space *spaces.Space, repo *Repository, artifact *Artifact, version string, filePath string) ([]byte, error) {
	// If the repository is an internal mirror, redirect the request to the mirrored repository
	if repo.InternalMirror != nil {
		mirroredSpace, err := s.spaces.Get(repo.InternalMirror.SpaceID)
		if err != nil {
			return nil, err
		}

		mirroredRepo, err := s.GetRepository(ctx, repo.InternalMirror.SpaceID, repo.InternalMirror.Repository)
		if err != nil {
			return nil, err
		}

		mirroredArtifact, err := s.GetArtifact(ctx, repo.InternalMirror.SpaceID, repo.InternalMirror.Repository, artifact.Group, artifact.ID)
		if err != nil {
			return nil, err
		}

		return s.GetJavadocFile(ctx, mirroredSpace, mirroredRepo, mirroredArtifact, version, filePath)
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%s", space.ID, repo.Name, artifact.Group, artifact.ID, version)

	if !s.javadocCache.IsJavadocCached(key) {
		javadocData, err := s.DownloadArtifactFile(ctx, space.ID, repo.Name, artifact.Group, artifact.ID, version, fmt.Sprintf("%s-%s-javadoc.jar", artifact.ID, version))
		if err != nil {
			return nil, err
		}

		files, err := s.javadocCache.CacheJavadoc(key, javadocData)
		if err != nil {
			return nil, err
		}
		data, exists := files[filePath]
		if !exists {
			return nil, javadoccache.ErrJavadocNotFound
		}
		return data, nil
	}

	return s.javadocCache.GetJavadocFile(key, filePath)
}
