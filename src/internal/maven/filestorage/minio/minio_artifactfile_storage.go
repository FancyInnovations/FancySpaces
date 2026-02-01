package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

const BucketName = "maven-artifact-files"

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

func (s *Storage) UploadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string, data []byte) error {
	objName := toObjName(spaceID, repoName, groupID, artifactID, version, fileName)
	r := bytes.NewReader(data)

	_, err := s.mio.PutObject(ctx, BucketName, objName, r, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DownloadArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) ([]byte, error) {
	objName := toObjName(spaceID, repoName, groupID, artifactID, version, fileName)
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

func (s *Storage) DeleteArtifactFile(ctx context.Context, spaceID, repoName, groupID, artifactID, version, fileName string) error {
	objName := toObjName(spaceID, repoName, groupID, artifactID, version, fileName)
	err := s.mio.RemoveObject(ctx, BucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func toObjName(spaceID, repoName, groupID, artifactID, version, fileName string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", spaceID, repoName, groupID, artifactID, version, fileName)
}
