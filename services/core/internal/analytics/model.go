package analytics

import (
	"time"
)

const AllVersionsID = "*"

type VersionDownload struct {
	SpaceID      string    `json:"space_id" ch:"space_id"`
	VersionID    string    `json:"version_id" ch:"version_id"`
	DownloadedAt time.Time `json:"downloaded_at" ch:"downloaded_at"`
	IPHash       string    `json:"ip_hash" ch:"ip_hash"`
	UserAgent    string    `json:"user_agent" ch:"user_agent"`
}

type MavenArtifactDownload struct {
	SpaceID        string    `json:"space_id" ch:"space_id"`
	RepositoryName string    `json:"repository_name" ch:"repository_name"`
	GroupID        string    `json:"group_id" ch:"group_id"`
	ArtifactID     string    `json:"artifact_id" ch:"artifact_id"`
	Version        string    `json:"version" ch:"version"`
	DownloadedAt   time.Time `json:"downloaded_at" ch:"downloaded_at"`
	IPHash         string    `json:"ip_hash" ch:"ip_hash"`
}
