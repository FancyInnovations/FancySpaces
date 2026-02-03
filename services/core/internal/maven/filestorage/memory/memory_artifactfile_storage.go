package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

var (
	ttl = 12 * time.Hour
)

type Storage struct {
	files *ristretto.Cache[string, []byte]
}

func NewStorage() *Storage {
	files, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 50 * 10,           // x10 of expected number of elements when full
		MaxCost:     512 * 1024 * 1024, // 512 MB
		BufferItems: 64,                // keep 64
	})
	if err != nil {
		panic(err)
	}

	return &Storage{
		files: files,
	}
}

func (s *Storage) UploadArtifactFile(_ context.Context, spaceID, repoName, groupID, artifactID, version, fileName string, data []byte) error {
	key := spaceID + ":" + repoName + ":" + groupID + ":" + artifactID + ":" + version + ":" + fileName
	s.files.SetWithTTL(key, data, int64(len(data)), ttl)
	return nil
}

func (s *Storage) DownloadArtifactFile(_ context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) ([]byte, error) {
	key := spaceID + ":" + repoName + ":" + groupID + ":" + artifactID + ":" + version + ":" + fileName
	data, found := s.files.Get(key)
	if !found {
		return nil, fmt.Errorf("file not found")
	}

	return data, nil
}

func (s *Storage) DeleteArtifactFile(_ context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) error {
	key := spaceID + ":" + repoName + ":" + groupID + ":" + artifactID + ":" + version + ":" + fileName
	s.files.Del(key)
	return nil
}
