package local

import (
	"context"
	"os"

	"github.com/fancyinnovations/fancyspaces/core/internal/versions"
)

const basePath = "data/versions"

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Upload(_ context.Context, version *versions.Version, file *versions.VersionFile, data []byte) error {

	dirPath := basePath + "/" + version.SpaceID + "/" + version.ID
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	filePath := dirPath + "/" + file.Name
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Download(_ context.Context, spaceID, versionID, fileName string) ([]byte, error) {
	filePath := basePath + "/" + spaceID + "/" + versionID + "/" + fileName
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Storage) Delete(_ context.Context, spaceID, versionID, fileName string) error {
	filePath := basePath + "/" + spaceID + "/" + versionID + "/" + fileName
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
