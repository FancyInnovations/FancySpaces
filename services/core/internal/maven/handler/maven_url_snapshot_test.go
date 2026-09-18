package handler

import "testing"

// TestSnapshotURLParsing verifies that a version ending in -SNAPSHOT parses
// unchanged across group/artifact/version/filename, exactly like a release URL.
func TestSnapshotURLParsing(t *testing.T) {
	url := "/space1/repo1/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-SNAPSHOT.jar"

	got, err := GroupFromURL(url)
	if err != nil {
		t.Fatalf("GroupFromURL: unexpected error: %v", err)
	}
	if got != "com.example" {
		t.Errorf("GroupFromURL: expected com.example, got %s", got)
	}

	got, err = ArtifactFromURL(url)
	if err != nil {
		t.Fatalf("ArtifactFromURL: unexpected error: %v", err)
	}
	if got != "myapp" {
		t.Errorf("ArtifactFromURL: expected myapp, got %s", got)
	}

	got, err = VersionFromURL(url)
	if err != nil {
		t.Fatalf("VersionFromURL: unexpected error: %v", err)
	}
	if got != "1.0.0-SNAPSHOT" {
		t.Errorf("VersionFromURL: expected 1.0.0-SNAPSHOT, got %s", got)
	}

	got, err = FilenameFromURL(url)
	if err != nil {
		t.Fatalf("FilenameFromURL: unexpected error: %v", err)
	}
	if got != "myapp-1.0.0-SNAPSHOT.jar" {
		t.Errorf("FilenameFromURL: expected myapp-1.0.0-SNAPSHOT.jar, got %s", got)
	}
}

// TestSnapshotMetadataURL confirms a metadata URL for a snapshot artifact is
// still recognized as metadata and extracts group, artifact, version correctly.
func TestSnapshotMetadataURL(t *testing.T) {
	url := "/space1/repo1/com/example/myapp/1.0.0-SNAPSHOT/maven-metadata.xml"
	if !IsMetadataURL(url) {
		t.Errorf("expected snapshot metadata URL to be metadata, got %v", IsMetadataURL(url))
	}

	got, err := GroupFromURL(url)
	if err != nil {
		t.Fatalf("GroupFromURL: unexpected error: %v", err)
	}
	if got != "com.example" {
		t.Errorf("GroupFromURL: expected com.example, got %s", got)
	}

	got, err = ArtifactFromURL(url)
	if err != nil {
		t.Fatalf("ArtifactFromURL: unexpected error: %v", err)
	}
	if got != "myapp" {
		t.Errorf("ArtifactFromURL: expected myapp, got %s", got)
	}

	got, err = VersionFromURL(url)
	if err != nil {
		t.Fatalf("VersionFromURL: unexpected error: %v", err)
	}
	if got != "1.0.0-SNAPSHOT" {
		t.Errorf("VersionFromURL: expected 1.0.0-SNAPSHOT, got %s", got)
	}

	got, err = FilenameFromURL(url)
	if err != nil {
		t.Fatalf("FilenameFromURL: unexpected error: %v", err)
	}
	if got != "maven-metadata.xml" {
		t.Errorf("FilenameFromURL: expected maven-metadata.xml, got %s", got)
	}
}

func TestTimestampedSnapshotURLParsing(t *testing.T) {
	url := "/space1/repo1/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-20240101.120000-1.jar"

	got, err := GroupFromURL(url)
	if err != nil {
		t.Fatalf("GroupFromURL: unexpected error: %v", err)
	}
	if got != "com.example" {
		t.Errorf("GroupFromURL: expected com.example, got %s", got)
	}

	got, err = ArtifactFromURL(url)
	if err != nil {
		t.Fatalf("ArtifactFromURL: unexpected error: %v", err)
	}
	if got != "myapp" {
		t.Errorf("ArtifactFromURL: expected myapp, got %s", got)
	}

	got, err = VersionFromURL(url)
	if err != nil {
		t.Fatalf("VersionFromURL: unexpected error: %v", err)
	}
	if got != "1.0.0-SNAPSHOT" {
		t.Errorf("VersionFromURL: expected 1.0.0-SNAPSHOT, got %s", got)
	}

	got, err = FilenameFromURL(url)
	if err != nil {
		t.Fatalf("FilenameFromURL: unexpected error: %v", err)
	}
	if got != "myapp-1.0.0-20240101.120000-1.jar" {
		t.Errorf("FilenameFromURL: expected myapp-1.0.0-20240101.120000-1.jar, got %s", got)
	}
}
