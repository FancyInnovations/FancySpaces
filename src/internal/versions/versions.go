package versions

import (
	"context"
	"errors"
)

type DB interface {
	GetAll(ctx context.Context, spaceID string) ([]Version, error)
	GetByID(ctx context.Context, spaceID, versionID string) (*Version, error)
	GetByName(ctx context.Context, spaceID, versionNumber string) (*Version, error)
	Create(ctx context.Context, v *Version) error
	Update(ctx context.Context, spaceID, versionID string, v *Version) error
	Delete(ctx context.Context, spaceID, versionID string) error
	LogDownload(ctx context.Context, spaceID, versionID string) error
}

type FileStorage interface {
	Upload(ctx context.Context, version *Version, file *VersionFile, data []byte) error
	Download(ctx context.Context, spaceID, versionID, fileName string) ([]byte, error)
	Delete(ctx context.Context, spaceID, versionID, fileName string) error
}

type Store struct {
	db          DB
	fileStorage FileStorage
}

type Configuration struct {
	DB          DB
	FileStorage FileStorage
}

func New(cfg Configuration) *Store {
	return &Store{
		db:          cfg.DB,
		fileStorage: cfg.FileStorage,
	}
}

func (s *Store) GetAll(ctx context.Context, spaceID string) ([]Version, error) {
	return s.db.GetAll(ctx, spaceID)
}

func (s *Store) Get(ctx context.Context, spaceID, id string) (*Version, error) {
	v, err := s.db.GetByID(ctx, spaceID, id)
	if err != nil {
		if errors.Is(err, ErrVersionNotFound) {
			v, err = s.db.GetByName(ctx, spaceID, id)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
		return nil, err
	}

	return v, nil
}

func (s *Store) Create(ctx context.Context, v *Version) error {
	return s.db.Create(ctx, v)
}

func (s *Store) Update(ctx context.Context, spaceID, versionID string, v *Version) error {
	return s.db.Update(ctx, spaceID, versionID, v)
}

func (s *Store) Delete(ctx context.Context, spaceID, versionID string) error {
	ver, err := s.Get(ctx, spaceID, versionID)
	if err != nil {
		return err
	}

	for _, f := range ver.Files {
		if err := s.fileStorage.Delete(ctx, spaceID, versionID, f.Name); err != nil {
			return err
		}
	}

	return s.db.Delete(ctx, spaceID, versionID)
}

func (s *Store) UploadVersionFile(ctx context.Context, version *Version, fileName string, data []byte) error {
	verFile := &VersionFile{
		Name: fileName,
		URL:  "https://fancyspaces.net/api/v1/spaces/" + version.SpaceID + "/versions/" + version.ID + "/files/" + fileName,
		Size: int64(len(data)),
	}

	version.Files = append(version.Files, *verFile)
	if err := s.Update(ctx, version.SpaceID, version.ID, version); err != nil {
		return err
	}

	return s.fileStorage.Upload(ctx, version, verFile, data)
}

func (s *Store) DownloadVersionFile(ctx context.Context, spaceID, versionID, fileName string) ([]byte, error) {
	ver, err := s.Get(ctx, spaceID, versionID)
	if err != nil {
		return nil, err
	}

	// Check if the file exists in the version
	found := false
	for _, f := range ver.Files {
		if f.Name == fileName {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrVersionNotFound
	}

	if err := s.db.LogDownload(ctx, spaceID, versionID); err != nil {
		return nil, err
	}

	return s.fileStorage.Download(ctx, spaceID, versionID, fileName)
}
