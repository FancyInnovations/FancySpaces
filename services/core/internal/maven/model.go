package maven

import (
	"encoding/xml"
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
	latest := ""
	release := ""
	if len(a.Versions) > 0 {
		latest = a.Versions[len(a.Versions)-1].Version
		release = a.Versions[len(a.Versions)-1].Version
	}

	versions := make([]string, len(a.Versions))
	for i, v := range a.Versions {
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

func (f *ArtifactVersion) GetFile(fileName string) *ArtifactVersionFile {
	for _, file := range f.Files {
		if file.Name == fileName {
			return file
		}
	}
	return nil
}
