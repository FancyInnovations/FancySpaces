package versions

import (
	"errors"
)

type DB interface {
	GetAll(spaceID string) ([]Version, error)
	GetByID(spaceID, versionID string) (*Version, error)
	GetByName(spaceID, versionNumber string) (*Version, error)
	Create(v *Version) error
	Update(spaceID, versionID string, v *Version) error
	Delete(spaceID, versionID string) error
	LogDownload(spaceID, versionID string) error
}

type FileStorage interface {
	Upload(version *Version, file *VersionFile, data []byte) error
	Download(spaceID, versionID, fileName string) ([]byte, error)
	Delete(spaceID, versionID, fileName string) error
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

func (s *Store) GetAll(spaceID string) ([]Version, error) {
	return s.db.GetAll(spaceID)
}

func (s *Store) Get(spaceID, id string) (*Version, error) {
	v, err := s.db.GetByID(spaceID, id)
	if err != nil {
		if errors.Is(err, ErrVersionNotFound) {
			v, err = s.db.GetByName(spaceID, id)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
		return nil, err
	}

	return v, nil
}

func (s *Store) Create(v *Version) error {
	return s.db.Create(v)
}

func (s *Store) Update(spaceID, versionID string, v *Version) error {
	return s.db.Update(spaceID, versionID, v)
}

func (s *Store) Delete(spaceID, versionID string) error {
	ver, err := s.Get(spaceID, versionID)
	if err != nil {
		return err
	}

	for _, f := range ver.Files {
		if err := s.fileStorage.Delete(spaceID, versionID, f.Name); err != nil {
			return err
		}
	}

	return s.db.Delete(spaceID, versionID)
}

func (s *Store) UploadVersionFile(version *Version, fileName string, data []byte) error {
	verFile := &VersionFile{
		Name: fileName,
		URL:  "https://fancyspaces/spaces/" + version.SpaceID + "/versions/" + version.ID + "/files/" + fileName,
		Size: int64(len(data)),
	}

	return s.fileStorage.Upload(version, verFile, data)
}

func (s *Store) DownloadVersionFile(spaceID, versionID, fileName string) ([]byte, error) {
	ver, err := s.Get(spaceID, versionID)
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

	if err := s.db.LogDownload(spaceID, versionID); err != nil {
		return nil, err
	}

	return s.fileStorage.Download(spaceID, versionID, fileName)
}
