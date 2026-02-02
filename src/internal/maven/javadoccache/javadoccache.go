package javadoccache

import (
	"archive/zip"
	"bytes"
	"io"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

var (
	ttl = 60 * time.Minute
)

type Service struct {
	fileCache *ristretto.Cache[string, map[string][]byte]
}

func NewService() *Service {
	fileCache, err := ristretto.NewCache(&ristretto.Config[string, map[string][]byte]{
		NumCounters: 50 * 10,           // x10 of expected number of elements when full
		MaxCost:     512 * 1024 * 1024, // 512 MB
		BufferItems: 64,                // keep 64
	})
	if err != nil {
		panic(err)
	}

	return &Service{
		fileCache: fileCache,
	}
}

func (s *Service) CacheJavadoc(key string, javadocZipData []byte) (map[string][]byte, error) {
	files, err := unzipToMap(javadocZipData)
	if err != nil {
		return nil, err
	}

	s.fileCache.SetWithTTL(key, files, int64(len(javadocZipData)), ttl)

	return files, nil
}

func (s *Service) GetJavadocFile(key, filePath string) ([]byte, error) {
	files, found := s.fileCache.Get(key)
	if !found {
		return nil, ErrJavadocNotFound
	}

	data, exists := files[filePath]
	if !exists {
		return nil, ErrJavadocNotFound
	}

	return data, nil
}

func (s *Service) IsJavadocCached(key string) bool {
	_, found := s.fileCache.Get(key)
	return found
}

func unzipToMap(zipData []byte) (map[string][]byte, error) {
	files := make(map[string][]byte)

	r := bytes.NewReader(zipData)
	zr, err := zip.NewReader(r, int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}

		files[f.Name] = data
	}

	return files, nil
}
