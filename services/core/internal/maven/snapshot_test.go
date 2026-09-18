package maven_test

import (
	"testing"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
)

func av(version string, at time.Time, files ...*maven.ArtifactVersionFile) *maven.ArtifactVersion {
	return &maven.ArtifactVersion{
		Version:     version,
		PublishedAt: at,
		Files:       files,
	}
}

func file(name string, size int64) *maven.ArtifactVersionFile {
	return &maven.ArtifactVersionFile{
		Name: name,
		Size: size,
		URL:  "https://maven.fancyspaces.net/space/repo/group/artifact/" + name,
	}
}

func newArtifact(v ...string) *maven.Artifact {
	a := &maven.Artifact{
		Group: "com.example",
		ID:    "myapp",
	}
	for _, ver := range v {
		a.Versions = append(a.Versions, av(ver, time.Now()))
	}
	return a
}

// GetVersion: resolves "latest", "release", and exact versions.
func TestGetVersion_Resolution(t *testing.T) {
	a := newArtifact("1.0.0", "1.1.0-SNAPSHOT")

	if got := a.GetVersion("latest").Version; got != "1.1.0-SNAPSHOT" {
		t.Fatalf("expected latest to resolve to newest version 1.1.0-SNAPSHOT, got %q", got)
	}
	if got := a.GetVersion("release").Version; got != "1.0.0" {
		t.Fatalf("expected release to resolve to latest release 1.0.0, got %q", got)
	}
	if got := a.GetVersion("1.0.0").Version; got != "1.0.0" {
		t.Fatalf("expected exact version 1.0.0, got %q", got)
	}
	if got := a.GetVersion("1.1.0-SNAPSHOT").Version; got != "1.1.0-SNAPSHOT" {
		t.Fatalf("expected exact snapshot 1.1.0-SNAPSHOT, got %q", got)
	}
}

// GetVersion: snapshot-only => latest is the snapshot, release is nil.
func TestGetVersion_SnapshotOnly(t *testing.T) {
	a := newArtifact("1.0.0-SNAPSHOT", "2.0.0-SNAPSHOT")
	if got := a.GetVersion("latest").Version; got != "2.0.0-SNAPSHOT" {
		t.Fatalf("expected latest to resolve to newest snapshot 2.0.0-SNAPSHOT, got %q", got)
	}
	if got := a.GetVersion("release"); got != nil {
		t.Fatalf("expected nil for release on snapshot-only artifact, got %v", got)
	}
}

// GetVersion: empty artifact => returns nil.
func TestGetVersion_Empty(t *testing.T) {
	a := &maven.Artifact{}
	if got := a.GetVersion("latest"); got != nil {
		t.Fatalf("expected nil for empty artifact latest, got %v", got)
	}
	if got := a.GetVersion("release"); got != nil {
		t.Fatalf("expected nil for empty artifact release, got %v", got)
	}
	if got := a.GetVersion("1.0.0"); got != nil {
		t.Fatalf("expected nil for missing exact version, got %v", got)
	}
}

// ToMetadataXML: release + snapshot => latest is newest version, release is newest release.
func TestToMetadataXML_ReleaseAndSnapshot(t *testing.T) {
	a := newArtifact("1.0.0", "1.1.0-SNAPSHOT")
	md := a.ToMetadataXML()
	if md.Versioning == nil {
		t.Fatalf("expected versioning in metadata")
	}
	if md.Versioning.Latest != "1.1.0-SNAPSHOT" {
		t.Fatalf("expected metadata latest 1.1.0-SNAPSHOT, got %q", md.Versioning.Latest)
	}
	if md.Versioning.Release != "1.0.0" {
		t.Fatalf("expected metadata release 1.0.0, got %q", md.Versioning.Release)
	}
	if len(md.Versioning.Versions) != 2 || md.Versioning.Versions[0] != "1.0.0" || md.Versioning.Versions[1] != "1.1.0-SNAPSHOT" {
		t.Fatalf("expected versions [1.0.0, 1.1.0-SNAPSHOT], got %v", md.Versioning.Versions)
	}
}

// ToMetadataXML: snapshot-only => latest is snapshot, release is empty.
func TestToMetadataXML_SnapshotOnly(t *testing.T) {
	a := newArtifact("2.0.0-SNAPSHOT")
	md := a.ToMetadataXML()
	if md.Versioning == nil {
		t.Fatalf("expected versioning in metadata")
	}
	if md.Versioning.Latest != "2.0.0-SNAPSHOT" {
		t.Fatalf("expected metadata latest 2.0.0-SNAPSHOT, got %q", md.Versioning.Latest)
	}
	if md.Versioning.Release != "" {
		t.Fatalf("expected metadata release to be empty for snapshot-only, got %q", md.Versioning.Release)
	}
	if len(md.Versioning.Versions) != 1 || md.Versioning.Versions[0] != "2.0.0-SNAPSHOT" {
		t.Fatalf("expected versions [2.0.0-SNAPSHOT], got %v", md.Versioning.Versions)
	}
}

// ToMetadataXML: release-only => latest and release are both the latest release.
func TestToMetadataXML_ReleaseOnly(t *testing.T) {
	a := newArtifact("1.0.0", "1.1.0", "1.2.0")
	md := a.ToMetadataXML()
	if md.Versioning == nil {
		t.Fatalf("expected versioning in metadata")
	}
	if md.Versioning.Latest != "1.2.0" {
		t.Fatalf("expected metadata latest 1.2.0, got %q", md.Versioning.Latest)
	}
	if md.Versioning.Release != "1.2.0" {
		t.Fatalf("expected metadata release 1.2.0, got %q", md.Versioning.Release)
	}
}

// ToVersionMetadataXML: generates version-level metadata.
func TestToVersionMetadataXML(t *testing.T) {
	pomFile := file("myapp-1.0.0-20240101.120000-1.pom", 1024)
	jarFile := file("myapp-1.0.0-20240101.120000-1.jar", 20480)
	sourcesFile := file("myapp-1.0.0-20240101.120000-1-sources.jar", 10240)

	version := av("1.0.0-SNAPSHOT", time.Now(), pomFile, jarFile, sourcesFile)
	a := &maven.Artifact{
		Group:    "com.example",
		ID:       "myapp",
		Versions: []*maven.ArtifactVersion{version},
	}

	md := a.ToVersionMetadataXML("1.0.0-SNAPSHOT")
	if md.GroupID != "com.example" {
		t.Fatalf("expected group com.example, got %q", md.GroupID)
	}
	if md.ArtifactID != "myapp" {
		t.Fatalf("expected artifact myapp, got %q", md.ArtifactID)
	}
	if md.Version != "1.0.0-SNAPSHOT" {
		t.Fatalf("expected version 1.0.0-SNAPSHOT, got %q", md.Version)
	}
	if md.Versioning == nil || md.Versioning.Snapshot == nil {
		t.Fatalf("expected versioning and snapshot in version metadata")
	}
	if md.Versioning.Snapshot.Timestamp != "20240101.120000" {
		t.Fatalf("expected timestamp 20240101.120000, got %s", md.Versioning.Snapshot.Timestamp)
	}
	if md.Versioning.Snapshot.BuildNumber != 1 {
		t.Fatalf("expected buildNumber 1, got %d", md.Versioning.Snapshot.BuildNumber)
	}
	if len(md.Versioning.SnapshotVersions) != 3 {
		t.Fatalf("expected 3 snapshotVersions, got %d", len(md.Versioning.SnapshotVersions))
	}
}

// Snapshot files: storing multiple files for a SNAPSHOT version preserves all files.
func TestArtifactVersion_SnapshotFilesPersistence(t *testing.T) {
	pomFile := file("myapp-1.0.0-SNAPSHOT.pom", 1024)
	jarFile := file("myapp-1.0.0-SNAPSHOT.jar", 20480)
	sourcesFile := file("myapp-1.0.0-SNAPSHOT-sources.jar", 10240)
	sha1File := file("myapp-1.0.0-SNAPSHOT.jar.sha1", 40)

	version := av("1.0.0-SNAPSHOT", time.Now(), pomFile, jarFile, sourcesFile, sha1File)
	a := &maven.Artifact{
		Group:    "com.example",
		ID:       "myapp",
		Versions: []*maven.ArtifactVersion{version},
	}

	v := a.GetVersion("1.0.0-SNAPSHOT")
	if v == nil {
		t.Fatalf("expected version 1.0.0-SNAPSHOT to exist")
	}
	if len(v.Files) != 4 {
		t.Fatalf("expected 4 files in snapshot version, got %d", len(v.Files))
	}
	if v.GetFile("myapp-1.0.0-SNAPSHOT.pom") == nil {
		t.Fatalf("expected pom file to be found")
	}
	if v.GetFile("myapp-1.0.0-SNAPSHOT.jar") == nil {
		t.Fatalf("expected jar file to be found")
	}
	if v.GetFile("myapp-1.0.0-SNAPSHOT-sources.jar") == nil {
		t.Fatalf("expected sources file to be found")
	}
	if v.GetFile("myapp-1.0.0-SNAPSHOT.jar.sha1") == nil {
		t.Fatalf("expected sha1 file to be found")
	}
}

func TestArtifactVersion_SnapshotResolutionFallback(t *testing.T) {
	pomFile := file("myapp-1.0.0-20240101.120000-1.pom", 1024)
	jarFile := file("myapp-1.0.0-20240101.120000-1.jar", 20480)
	sourcesFile := file("myapp-1.0.0-20240101.120000-1-sources.jar", 10240)
	sha1File := file("myapp-1.0.0-20240101.120000-1.jar.sha1", 40)

	version := av("1.0.0-SNAPSHOT", time.Now(), pomFile, jarFile, sourcesFile, sha1File)
	a := &maven.Artifact{
		Group:    "com.example",
		ID:       "myapp",
		Versions: []*maven.ArtifactVersion{version},
	}

	v := a.GetVersion("1.0.0-SNAPSHOT")
	if v == nil {
		t.Fatalf("expected version 1.0.0-SNAPSHOT to exist")
	}

	// Requesting generic SNAPSHOT name resolves to timestamped build
	if got := v.GetFile("myapp-1.0.0-SNAPSHOT.jar"); got == nil || got.Name != "myapp-1.0.0-20240101.120000-1.jar" {
		t.Fatalf("expected resolution to myapp-1.0.0-20240101.120000-1.jar, got %v", got)
	}
	if got := v.GetFile("myapp-1.0.0-SNAPSHOT.pom"); got == nil || got.Name != "myapp-1.0.0-20240101.120000-1.pom" {
		t.Fatalf("expected resolution to myapp-1.0.0-20240101.120000-1.pom, got %v", got)
	}
	if got := v.GetFile("myapp-1.0.0-SNAPSHOT-sources.jar"); got == nil || got.Name != "myapp-1.0.0-20240101.120000-1-sources.jar" {
		t.Fatalf("expected resolution to myapp-1.0.0-20240101.120000-1-sources.jar, got %v", got)
	}
	if got := v.GetFile("myapp-1.0.0-SNAPSHOT.jar.sha1"); got == nil || got.Name != "myapp-1.0.0-20240101.120000-1.jar.sha1" {
		t.Fatalf("expected resolution to myapp-1.0.0-20240101.120000-1.jar.sha1, got %v", got)
	}
}

// Multiple snapshot versions coexist without deleting each other.
func TestMultipleSnapshotVersionsCoexist(t *testing.T) {
	v1 := av("1.0.0-SNAPSHOT", time.Now(), file("myapp-1.0.0-SNAPSHOT.pom", 500))
	v2 := av("1.1.0-SNAPSHOT", time.Now(), file("myapp-1.1.0-SNAPSHOT.pom", 600))
	v3 := av("2.0.0-SNAPSHOT", time.Now(), file("myapp-2.0.0-SNAPSHOT.pom", 700))

	a := &maven.Artifact{
		Group:    "com.example",
		ID:       "myapp",
		Versions: []*maven.ArtifactVersion{v1, v2, v3},
	}

	if len(a.Versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(a.Versions))
	}
	if a.GetVersion("1.0.0-SNAPSHOT") == nil {
		t.Fatalf("expected 1.0.0-SNAPSHOT to exist")
	}
	if a.GetVersion("1.1.0-SNAPSHOT") == nil {
		t.Fatalf("expected 1.1.0-SNAPSHOT to exist")
	}
	if a.GetVersion("2.0.0-SNAPSHOT") == nil {
		t.Fatalf("expected 2.0.0-SNAPSHOT to exist")
	}
}
