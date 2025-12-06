package analytics

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

type DB interface {
	GetDownloadCountForSpace(ctx context.Context, spaceID string) (int64, error)
	GetDownloadCountForVersion(ctx context.Context, spaceID, versionID string) (int64, error)
	StoreVersionDownloads(ctx context.Context, records []VersionDownload) error
}

type Cache interface {
	GetDownloadCountForVersion(spaceID, versionID string) int64
	SetDownloadCountForVersion(spaceID, versionID string, count int64)
}

type Store struct {
	db DB
	c  Cache
}

type Configuration struct {
	DB    DB
	Cache Cache
}

func New(cfg Configuration) *Store {
	return &Store{
		db: cfg.DB,
		c:  cfg.Cache,
	}
}

func (s *Store) GetDownloadCountForVersion(ctx context.Context, spaceID, versionID string) (int64, error) {
	if count := s.c.GetDownloadCountForVersion(spaceID, versionID); count >= 0 {
		return count, nil
	}

	count, err := s.db.GetDownloadCountForVersion(ctx, spaceID, versionID)
	if err != nil {
		return -1, err
	}

	s.c.SetDownloadCountForVersion(spaceID, versionID, count)

	return count, nil
}

func (s *Store) LogDownloadForVersion(ctx context.Context, spaceID, versionID string, r *http.Request) error {
	ip := getIP(r)
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

func getIP(r *http.Request) string {
	ip := "unknown"
	if r.RemoteAddr != "" {
		ip = r.RemoteAddr
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		ip = xri
	}
	return ip
}

func hashIP(ip string) string {
	h := sha256.New()
	h.Write([]byte(ip))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
