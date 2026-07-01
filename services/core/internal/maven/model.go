package maven

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	SpaceID        string          `json:"space_id" bson:"space_id"`
	Name           string          `json:"name" bson:"name"`
	Public         bool            `json:"public" bson:"public"`
	CreatedAt      time.Time       `json:"created_at" bson:"created_at"`
	InternalMirror *InternalMirror `json:"internal_mirror,omitempty" bson:"internal_mirror,omitempty"`
}

type InternalMirror struct {
	SpaceID    string `json:"space_id" bson:"space_id"`
	Repository string `json:"repository" bson:"repository"`
}

type Artifact struct {
	SpaceID    string             `json:"space_id" bson:"space_id"`
	Repository string             `json:"repository" bson:"repository"`
	Group      string             `json:"group" bson:"group"`
	ID         string             `json:"id" bson:"id"`
	Versions   []*ArtifactVersion `json:"versions" bson:"versions"`
}

type ArtifactVersion struct {
	Version     string                 `json:"version" bson:"version"`
	PublishedAt time.Time              `json:"published_at" bson:"published_at"`
	Files       []*ArtifactVersionFile `json:"files" bson:"files"`
}

type ArtifactVersionFile struct {
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	URL  string `json:"url" bson:"url"`
}

type MetadataXML struct {
	XMLName    xml.Name `xml:"metadata"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Versions   []string `xml:"versioning>versions>version"`
	Latest     string   `xml:"versioning>latest"`
	Release    string   `xml:"versioning>release"`
}

func (a *Artifact) GetVersion(version string) *ArtifactVersion {
	// "latest" returns the most recent stored version
	if version == "latest" && len(a.Versions) > 0 {
		return a.Versions[len(a.Versions)-1]
	}

	// direct match
	for _, v := range a.Versions {
		if v.Version == version {
			return v
		}
	}

	// support SNAPSHOT resolution: when client requests X-SNAPSHOT, find the
	// latest stored timestamped snapshot like X-20230701.123456-1
	if strings.HasSuffix(version, "-SNAPSHOT") {
		base := strings.TrimSuffix(version, "-SNAPSHOT")
		var found *ArtifactVersion
		for _, v := range a.Versions {
			if strings.HasPrefix(v.Version, base+"-") {
				// keep iterating so the last matching (most recent) is returned
				found = v
			}
		}
		if found != nil {
			return found
		}
	}

	return nil
}

func (a *Artifact) ToMetadataXML() MetadataXML {
	// sort versions deterministically by PublishedAt then by Version
	sorted := make([]*ArtifactVersion, len(a.Versions))
	copy(sorted, a.Versions)
	sort.Slice(sorted, func(i, j int) bool {
		if !sorted[i].PublishedAt.Equal(sorted[j].PublishedAt) {
			return sorted[i].PublishedAt.Before(sorted[j].PublishedAt)
		}
		return sorted[i].Version < sorted[j].Version
	})

	latest := ""
	release := ""
	if len(sorted) > 0 {
		latest = sorted[len(sorted)-1].Version
		// find last non-SNAPSHOT as release
		for k := len(sorted) - 1; k >= 0; k-- {
			if !strings.Contains(sorted[k].Version, "SNAPSHOT") {
				release = sorted[k].Version
				break
			}
		}
		// fallback to latest if no non-snapshot found
		if release == "" {
			release = latest
		}
	}

	versions := make([]string, len(sorted))
	for i, v := range sorted {
		versions[i] = v.Version
	}

	return MetadataXML{
		GroupID:    a.Group,
		ArtifactID: a.ID,
		Version:    latest,
		Versions:   versions,
		Latest:     latest,
		Release:    release,
	}
}

type SnapshotVersion struct {
	Extension  string `xml:"extension"`
	Classifier string `xml:"classifier,omitempty"`
	Value      string `xml:"value"`
	Updated    string `xml:"updated"`
}

type VersionMetadataXML struct {
	XMLName    xml.Name `xml:"metadata"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Versioning struct {
		Snapshot         *struct {
			Timestamp   string `xml:"timestamp"`
			BuildNumber int    `xml:"buildNumber"`
		} `xml:"snapshot,omitempty"`
		LastUpdated      string            `xml:"lastUpdated"`
		SnapshotVersions []SnapshotVersion `xml:"snapshotVersions>snapshotVersion"`
	} `xml:"versioning"`
}

func (a *Artifact) ToSnapshotMetadataXML(snapshotVersion string) VersionMetadataXML {
	base := strings.TrimSuffix(snapshotVersion, "-SNAPSHOT")
	// collect timestamped snapshots that match base
	type sv struct {
		version     string
		publishedAt time.Time
		files       []*ArtifactVersionFile
	}
	var snaps []sv
	for _, v := range a.Versions {
		if strings.HasPrefix(v.Version, base+"-") {
			snaps = append(snaps, sv{version: v.Version, publishedAt: v.PublishedAt, files: v.Files})
		}
	}

	// sort snaps by PublishedAt ascending
	sort.Slice(snaps, func(i, j int) bool {
		if !snaps[i].publishedAt.Equal(snaps[j].publishedAt) {
			return snaps[i].publishedAt.Before(snaps[j].publishedAt)
		}
		return snaps[i].version < snaps[j].version
	})

	var latest time.Time
	latestTimestamp := ""
	latestBuildNumber := 0
	if len(snaps) > 0 {
		latest = snaps[len(snaps)-1].publishedAt
		// parse timestamp/build from last snap version
		suffix := strings.TrimPrefix(snaps[len(snaps)-1].version, base+"-")
		parts := strings.Split(suffix, "-")
		if len(parts) >= 2 {
			latestTimestamp = parts[0]
			latestBuildNumber, _ = strconv.Atoi(parts[1])
		}
	}

	// build deduped snapshotVersions keyed by extension|classifier
	dedupe := map[string]SnapshotVersion{}
	for _, s := range snaps {
		updated := s.publishedAt.UTC().Format("20060102150405")
		for _, f := range s.files {
			ext := ""
			if idx := strings.LastIndex(f.Name, "."); idx != -1 {
				ext = f.Name[idx+1:]
			}

			// determine classifier: what's between artifactId-version- and .ext
			cls := ""
			prefix := a.ID + "-" + s.version + "-"
			if strings.HasPrefix(f.Name, prefix) {
				clsWithExt := f.Name[len(prefix):]
				if dot := strings.LastIndex(clsWithExt, "."); dot != -1 {
					cls = clsWithExt[:dot]
				} else {
					cls = clsWithExt
				}
			}

			key := ext + "|" + cls
			// keep the latest entry for each key
			existing, ok := dedupe[key]
			if !ok || existing.Updated < updated {
				dedupe[key] = SnapshotVersion{
					Extension:  ext,
					Classifier: cls,
					Value:      s.version,
					Updated:    updated,
				}
			}
		}
	}

	// convert dedupe map to slice and sort for deterministic output
	snapshotVersions := make([]SnapshotVersion, 0, len(dedupe))
	for _, v := range dedupe {
		snapshotVersions = append(snapshotVersions, v)
	}
	sort.Slice(snapshotVersions, func(i, j int) bool {
		if snapshotVersions[i].Extension != snapshotVersions[j].Extension {
			return snapshotVersions[i].Extension < snapshotVersions[j].Extension
		}
		return snapshotVersions[i].Classifier < snapshotVersions[j].Classifier
	})

	vm := VersionMetadataXML{GroupID: a.Group, ArtifactID: a.ID, Version: snapshotVersion}
	if latestTimestamp != "" {
		vm.Versioning.Snapshot = &struct {
			Timestamp   string `xml:"timestamp"`
			BuildNumber int    `xml:"buildNumber"`
		}{latestTimestamp, latestBuildNumber}
	}
	if !latest.IsZero() {
		vm.Versioning.LastUpdated = latest.UTC().Format("20060102150405")
	}
	vm.Versioning.SnapshotVersions = snapshotVersions

	return vm
}

func (f *ArtifactVersion) GetFile(fileName string) *ArtifactVersionFile {
	for _, file := range f.Files {
		if file.Name == fileName {
			return file
		}
	}
	return nil
}
