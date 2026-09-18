package maven

import (
	"encoding/xml"
	"fmt"
	"regexp"
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
	XMLName      xml.Name       `xml:"metadata"`
	ModelVersion string         `xml:"modelVersion,attr,omitempty"`
	GroupID      string         `xml:"groupId"`
	ArtifactID   string         `xml:"artifactId"`
	Version      string         `xml:"version,omitempty"`
	Versioning   *VersioningXML `xml:"versioning,omitempty"`
}

type VersioningXML struct {
	Latest           string               `xml:"latest,omitempty"`
	Release          string               `xml:"release,omitempty"`
	Snapshot         *SnapshotXML         `xml:"snapshot,omitempty"`
	Versions         []string             `xml:"versions>version,omitempty"`
	LastUpdated      string               `xml:"lastUpdated,omitempty"`
	SnapshotVersions []SnapshotVersionXML `xml:"snapshotVersions>snapshotVersion,omitempty"`
}

type SnapshotXML struct {
	Timestamp   string `xml:"timestamp,omitempty"`
	BuildNumber int    `xml:"buildNumber,omitempty"`
	LocalCopy   bool   `xml:"localCopy,omitempty"`
}

type SnapshotVersionXML struct {
	Classifier string `xml:"classifier,omitempty"`
	Extension  string `xml:"extension"`
	Value      string `xml:"value"`
	Updated    string `xml:"updated"`
}

var (
	timestampedSnapshotRe = regexp.MustCompile(`^(?:(.+)-)?(\d{8}\.\d{6})-(\d+)(?:-([a-zA-Z0-9_-]+))?\.([a-zA-Z0-9]+(?:\.(?:sha1|md5|sha256|sha512|asc))?)$`)
	genericSnapshotRe     = regexp.MustCompile(`^(?:(.+)-)?SNAPSHOT(?:-([a-zA-Z0-9_-]+))?\.([a-zA-Z0-9]+(?:\.(?:sha1|md5|sha256|sha512|asc))?)$`)
)

func isChecksumOrSignatureExt(ext string) bool {
	return strings.HasSuffix(ext, ".sha1") ||
		strings.HasSuffix(ext, ".md5") ||
		strings.HasSuffix(ext, ".sha256") ||
		strings.HasSuffix(ext, ".sha512") ||
		strings.HasSuffix(ext, ".asc")
}

// isSnapshot reports whether the version string is a Maven snapshot marker.
func isSnapshot(version string) bool {
	return strings.HasSuffix(version, "-SNAPSHOT")
}

// resolveLatest returns the newest version overall (the last element).
func (a *Artifact) resolveLatest() *ArtifactVersion {
	if len(a.Versions) == 0 {
		return nil
	}
	return a.Versions[len(a.Versions)-1]
}

// releaseVersion returns the newest non-snapshot release version, or nil when the
// artifact has no release (snapshot-only).
func (a *Artifact) releaseVersion() *ArtifactVersion {
	for i := len(a.Versions) - 1; i >= 0; i-- {
		if !isSnapshot(a.Versions[i].Version) {
			return a.Versions[i]
		}
	}
	return nil
}

func (a *Artifact) GetVersion(version string) *ArtifactVersion {
	if version == "latest" {
		return a.resolveLatest()
	}
	if version == "release" {
		return a.releaseVersion()
	}

	for _, v := range a.Versions {
		if v.Version == version {
			return v
		}
	}
	return nil
}

func (a *Artifact) ToMetadataXML() MetadataXML {
	latest := ""
	release := ""
	if len(a.Versions) > 0 {
		if r := a.resolveLatest(); r != nil {
			latest = r.Version
		}
		if rv := a.releaseVersion(); rv != nil {
			release = rv.Version
		}
	}

	versions := make([]string, len(a.Versions))
	var lastUpdatedTime time.Time
	for i, v := range a.Versions {
		versions[i] = v.Version
		if v.PublishedAt.After(lastUpdatedTime) {
			lastUpdatedTime = v.PublishedAt
		}
	}

	var lastUpdated string
	if !lastUpdatedTime.IsZero() {
		lastUpdated = lastUpdatedTime.UTC().Format("20060102150405")
	}

	return MetadataXML{
		GroupID:    a.Group,
		ArtifactID: a.ID,
		Versioning: &VersioningXML{
			Latest:      latest,
			Release:     release,
			Versions:    versions,
			LastUpdated: lastUpdated,
		},
	}
}

func (a *Artifact) ToVersionMetadataXML(version string) MetadataXML {
	v := a.GetVersion(version)
	if v == nil || !isSnapshot(version) {
		lastUpdated := ""
		if v != nil && !v.PublishedAt.IsZero() {
			lastUpdated = v.PublishedAt.UTC().Format("20060102150405")
		}
		return MetadataXML{
			GroupID:    a.Group,
			ArtifactID: a.ID,
			Version:    version,
			Versioning: &VersioningXML{
				LastUpdated: lastUpdated,
			},
		}
	}

	baseVersion := strings.TrimSuffix(version, "-SNAPSHOT")
	var (
		latestTimestamp   string
		latestBuildNumber int
		snapshotVersions  []SnapshotVersionXML
		seenSnapshots     = make(map[string]bool)
	)

	lastUpdated := v.PublishedAt.UTC().Format("20060102150405")

	for _, file := range v.Files {
		if m := timestampedSnapshotRe.FindStringSubmatch(file.Name); len(m) > 0 {
			ts := m[2]
			bn, _ := strconv.Atoi(m[3])
			classifier := m[4]
			ext := m[5]

			if ts > latestTimestamp || (ts == latestTimestamp && bn > latestBuildNumber) {
				latestTimestamp = ts
				latestBuildNumber = bn
			}

			if !isChecksumOrSignatureExt(ext) {
				val := fmt.Sprintf("%s-%s-%d", baseVersion, ts, bn)
				key := classifier + ":" + ext + ":" + val
				if !seenSnapshots[key] {
					seenSnapshots[key] = true
					snapshotVersions = append(snapshotVersions, SnapshotVersionXML{
						Classifier: classifier,
						Extension:  ext,
						Value:      val,
						Updated:    strings.ReplaceAll(ts, ".", ""),
					})
				}
			}
		} else if m := genericSnapshotRe.FindStringSubmatch(file.Name); len(m) > 0 {
			classifier := m[2]
			ext := m[3]
			if !isChecksumOrSignatureExt(ext) {
				key := classifier + ":" + ext + ":" + version
				if !seenSnapshots[key] {
					seenSnapshots[key] = true
					snapshotVersions = append(snapshotVersions, SnapshotVersionXML{
						Classifier: classifier,
						Extension:  ext,
						Value:      version,
						Updated:    lastUpdated,
					})
				}
			}
		}
	}

	var snapshot *SnapshotXML
	if latestTimestamp != "" {
		snapshot = &SnapshotXML{
			Timestamp:   latestTimestamp,
			BuildNumber: latestBuildNumber,
		}
		lastUpdated = strings.ReplaceAll(latestTimestamp, ".", "")
	} else {
		snapshot = &SnapshotXML{
			LocalCopy: true,
		}
	}

	return MetadataXML{
		ModelVersion: "1.1.0",
		GroupID:      a.Group,
		ArtifactID:   a.ID,
		Version:      version,
		Versioning: &VersioningXML{
			Snapshot:         snapshot,
			LastUpdated:      lastUpdated,
			SnapshotVersions: snapshotVersions,
		},
	}
}

func (f *ArtifactVersion) GetFile(fileName string) *ArtifactVersionFile {
	for _, file := range f.Files {
		if file.Name == fileName {
			return file
		}
	}

	if !isSnapshot(f.Version) {
		return nil
	}

	// Snapshot resolution: match generic snapshot request (e.g. myapp-1.0.0-SNAPSHOT.jar)
	// to timestamped file (e.g. myapp-1.0.0-20240101.120000-1.jar), or vice versa.
	if m := genericSnapshotRe.FindStringSubmatch(fileName); len(m) > 0 {
		prefix := m[1]
		classifier := m[2]
		ext := m[3]

		var bestMatch *ArtifactVersionFile
		var bestTS string
		var bestBN int

		for _, file := range f.Files {
			if tm := timestampedSnapshotRe.FindStringSubmatch(file.Name); len(tm) > 0 {
				tPrefix := tm[1]
				tTS := tm[2]
				tBN, _ := strconv.Atoi(tm[3])
				tClassifier := tm[4]
				tExt := tm[5]

				if tPrefix == prefix && tClassifier == classifier && tExt == ext {
					if tTS > bestTS || (tTS == bestTS && tBN > bestBN) {
						bestTS = tTS
						bestBN = tBN
						bestMatch = file
					}
				}
			}
		}
		if bestMatch != nil {
			return bestMatch
		}
	} else if tm := timestampedSnapshotRe.FindStringSubmatch(fileName); len(tm) > 0 {
		prefix := tm[1]
		classifier := tm[4]
		ext := tm[5]

		for _, file := range f.Files {
			if gm := genericSnapshotRe.FindStringSubmatch(file.Name); len(gm) > 0 {
				gPrefix := gm[1]
				gClassifier := gm[2]
				gExt := gm[3]

				if gPrefix == prefix && gClassifier == classifier && gExt == ext {
					return file
				}
			}
		}
	}

	return nil
}
