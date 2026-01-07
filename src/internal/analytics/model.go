package analytics

import (
	"errors"
	"time"
)

type VersionDownload struct {
	SpaceID      string    `json:"space_id" ch:"space_id"`
	VersionID    string    `json:"version_id" ch:"version_id"`
	DownloadedAt time.Time `json:"downloaded_at" ch:"downloaded_at"`
	IPHash       string    `json:"ip_hash" ch:"ip_hash"`
	UserAgent    string    `json:"user_agent" ch:"user_agent"`
}

var (
	NotInCacheErr = errors.New("not in cache")
)
