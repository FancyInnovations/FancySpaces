package analytics

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

type DB interface {
	GetDownloadCountForSpace(ctx context.Context, spaceID string) (uint64, error)
	GetDownloadCountForVersion(ctx context.Context, spaceID, versionID string) (uint64, error)
	StoreVersionDownloads(ctx context.Context, records []VersionDownload) error
}

type Cache interface {
	GetDownloadCountForVersion(spaceID, versionID string) (error, uint64)
	SetDownloadCountForVersion(spaceID, versionID string, count uint64)
}

type Store struct {
	db    DB
	c     Cache
	getIP func(r *http.Request) string
}

type Configuration struct {
	DB    DB
	Cache Cache
	GetIP func(r *http.Request) string
}

func New(cfg Configuration) *Store {
	return &Store{
		db:    cfg.DB,
		c:     cfg.Cache,
		getIP: cfg.GetIP,
	}
}

func (s *Store) GetDownloadCountForSpace(ctx context.Context, spaceID string) (uint64, error) {
	err, count := s.c.GetDownloadCountForVersion(spaceID, "*")
	if err == nil {
		return count, nil
	}

	count, err = s.db.GetDownloadCountForSpace(ctx, spaceID)
	if err != nil {
		return 0, err
	}

	s.c.SetDownloadCountForVersion(spaceID, "*", count)

	return count, nil
}

func (s *Store) GetDownloadCountForVersion(ctx context.Context, spaceID, versionID string) (uint64, error) {
	err, count := s.c.GetDownloadCountForVersion(spaceID, versionID)
	if err == nil {
		return count, nil
	}

	count, err = s.db.GetDownloadCountForVersion(ctx, spaceID, versionID)
	if err != nil {
		return 0, err
	}

	s.c.SetDownloadCountForVersion(spaceID, versionID, count)

	return count, nil
}

func (s *Store) LogDownloadForVersion(ctx context.Context, spaceID, versionID string, r *http.Request) error {
	ip := s.getIP(r)
	if ip != "unknown" {
		ip = hashIP(ip)
	}

	userAgent := r.UserAgent()

	vd := VersionDownload{
		SpaceID:      spaceID,
		VersionID:    versionID,
		DownloadedAt: time.Now(),
		IPHash:       ip,
		UserAgent:    userAgent,
	}

	return s.db.StoreVersionDownloads(ctx, []VersionDownload{vd}) // TODO: batch inserts
}

func hashIP(ip string) string {
	h := sha256.New()
	h.Write([]byte(ip))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
