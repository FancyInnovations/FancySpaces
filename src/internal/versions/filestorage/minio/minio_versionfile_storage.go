package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/src/internal/versions"
	"github.com/minio/minio-go/v7"
)

const BucketName = "version-files"

type Storage struct {
	mio *minio.Client
}

func NewStorage(mio *minio.Client) *Storage {
	return &Storage{mio: mio}
}

func (s *Storage) Setup(ctx context.Context) error {
	err := s.mio.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := s.mio.BucketExists(ctx, BucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

func (s *Storage) Upload(ctx context.Context, version *versions.Version, file *versions.VersionFile, data []byte) error {
	objName := fmt.Sprintf("%s/%s/%s", version.SpaceID, version.ID, file.Name)
	r := bytes.NewReader(data)

	_, err := s.mio.PutObject(ctx, BucketName, objName, r, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Download(ctx context.Context, spaceID, versionID, fileName string) ([]byte, error) {
	objName := fmt.Sprintf("%s/%s/%s", spaceID, versionID, fileName)
	obj, err := s.mio.GetObject(ctx, BucketName, objName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Storage) Delete(ctx context.Context, spaceID, versionID, fileName string) error {
	objName := fmt.Sprintf("%s/%s/%s", spaceID, versionID, fileName)
	err := s.mio.RemoveObject(ctx, BucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
