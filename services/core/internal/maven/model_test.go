package maven

import (
	"testing"
	"time"
)

func TestArtifact_GetVersionAndGetFile(t *testing.T) {
	a := Artifact{
		Group: "com.example",
		ID:    "demo",
		Versions: []*ArtifactVersion{
			{Version: "1.0-20230701.123456-1", PublishedAt: time.Date(2023, 7, 1, 12, 34, 56, 0, time.UTC), Files: []*ArtifactVersionFile{{Name: "demo-1.0-20230701.123456-1.jar"}}},
			{Version: "1.0-20230702.000001-2", PublishedAt: time.Date(2023, 7, 2, 0, 0, 1, 0, time.UTC), Files: []*ArtifactVersionFile{{Name: "demo-1.0-20230702.000001-2.jar"}, {Name: "demo-1.0-20230702.000001-2-sources.jar"}}},
			{Version: "2.0", PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Files: []*ArtifactVersionFile{{Name: "demo-2.0.jar"}}},
		},
	}

	// latest
	v := a.GetVersion("latest")
	if v == nil || v.Version != "2.0" {
		t.Fatalf("expected latest version 2.0, got %v", v)
	}

	// direct match
	v = a.GetVersion("1.0-20230701.123456-1")
	if v == nil || v.Version != "1.0-20230701.123456-1" {
		t.Fatalf("expected direct version match, got %v", v)
	}

	// snapshot resolution
	v = a.GetVersion("1.0-SNAPSHOT")
	if v == nil || v.Version != "1.0-20230702.000001-2" {
		t.Fatalf("expected snapshot to resolve to latest timestamped snapshot, got %v", v)
	}

	// GetFile
	file := v.GetFile("demo-1.0-20230702.000001-2-sources.jar")
	if file == nil || file.Name != "demo-1.0-20230702.000001-2-sources.jar" {
		t.Fatalf("expected to find file, got %v", file)
	}
}

func TestArtifact_ToMetadataXML(t *testing.T) {
	a := Artifact{
		Group: "com.example",
		ID:    "demo",
		Versions: []*ArtifactVersion{
			{Version: "1.0-20230701.123456-1"},
			{Version: "1.0-20230702.000001-2"},
			{Version: "2.0"},
		},
	}

	m := a.ToMetadataXML()
	if m.Latest != "2.0" {
		t.Fatalf("expected latest 2.0, got %s", m.Latest)
	}
	if m.Release != "2.0" {
		t.Fatalf("expected release 2.0, got %s", m.Release)
	}
	if len(m.Versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(m.Versions))
	}
}

func TestArtifact_ToSnapshotMetadataXML(t *testing.T) {
	a := Artifact{
		Group: "com.example",
		ID:    "demo",
		Versions: []*ArtifactVersion{
			{Version: "1.0-20230701.123456-1", PublishedAt: time.Date(2023, 7, 1, 12, 34, 56, 0, time.UTC), Files: []*ArtifactVersionFile{{Name: "demo-1.0-20230701.123456-1.jar"}}},
			{Version: "1.0-20230702.000001-2", PublishedAt: time.Date(2023, 7, 2, 0, 0, 1, 0, time.UTC), Files: []*ArtifactVersionFile{{Name: "demo-1.0-20230702.000001-2.jar"}, {Name: "demo-1.0-20230702.000001-2-sources.jar"}}},
		},
	}

	vm := a.ToSnapshotMetadataXML("1.0-SNAPSHOT")
	if vm.Version != "1.0-SNAPSHOT" {
		t.Fatalf("expected version to be 1.0-SNAPSHOT, got %s", vm.Version)
	}
	if vm.Versioning.Snapshot == nil {
		t.Fatalf("expected snapshot element to be present")
	}
	if vm.Versioning.Snapshot.Timestamp != "20230702.000001" {
		t.Fatalf("expected timestamp 20230702.000001, got %s", vm.Versioning.Snapshot.Timestamp)
	}
	if vm.Versioning.Snapshot.BuildNumber != 2 {
		t.Fatalf("expected build number 2, got %d", vm.Versioning.Snapshot.BuildNumber)
	}
	if vm.Versioning.LastUpdated == "" {
		t.Fatalf("expected lastUpdated to be set")
	}
	if len(vm.Versioning.SnapshotVersions) != 3 {
		t.Fatalf("expected 3 snapshotVersions entries (2 files + 1 jar), got %d", len(vm.Versioning.SnapshotVersions))
	}
}
