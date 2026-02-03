package handler

import "testing"

func TestGroupFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "com.example",
		},
		{
			url:      "/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "org.apache.commons",
		},
		{
			url:      "/space3/repo3/net/sf/jopt-simple/5.0/jopt-simple-5.0.pom",
			expected: "net.sf",
		},
		{
			url:      "/space4/repo4/io/github/user/project/maven-metadata.xml",
			expected: "io.github.user.project",
		},
		{
			url:      "/space5/repo5/com/example/lib/maven-metadata.xml.sha1",
			expected: "com.example.lib",
		},
	}

	for _, tc := range tests {
		group, err := GroupFromURL(tc.url)
		if err != nil {
			t.Errorf("Unexpected error for URL %s: %v", tc.url, err)
		}

		if group != tc.expected {
			t.Errorf("For URL %s, expected group %s but got %s", tc.url, tc.expected, group)
		}
	}
}

func TestArtifactFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "myapp",
		},
		{
			url:      "/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "lang3",
		},
		{
			url:      "/space3/repo3/net/sf/jopt-simple/5.0/jopt-simple-5.0.pom",
			expected: "jopt-simple",
		},
		{
			url:      "/space4/repo4/io/github/user/project/maven-metadata.xml",
			expected: "project",
		},
		{
			url:      "/space5/repo5/com/example/lib/maven-metadata.xml.sha1",
			expected: "lib",
		},
	}

	for _, tc := range tests {
		artifact, err := ArtifactFromURL(tc.url)
		if err != nil {
			t.Errorf("Unexpected error for URL %s: %v", tc.url, err)
		}

		if artifact != tc.expected {
			t.Errorf("For URL %s, expected artifact %s but got %s", tc.url, tc.expected, artifact)
		}
	}
}

func TestVersionFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "1.0.0",
		},
		{
			url:      "/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "3.12.0",
		},
	}

	for _, tc := range tests {
		version, err := VersionFromURL(tc.url)
		if err != nil {
			t.Errorf("Unexpected error for URL %s: %v", tc.url, err)
		}

		if version != tc.expected {
			t.Errorf("For URL %s, expected version %s but got %s", tc.url, tc.expected, version)
		}
	}
}

func TestIsMetadataURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{
			url:      "/space4/repo4/io/github/user/project/maven-metadata.xml",
			expected: true,
		},
		{
			url:      "/space5/repo5/com/example/lib/maven-metadata.xml.sha1",
			expected: true,
		},
		{
			url:      "/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: false,
		},
		{
			url:      "/space1/repo1/com/example/lib/maven-metadata.xmlish",
			expected: false,
		},
		{
			url:      "maven-metadata.xml",
			expected: true,
		},
		{
			url:      "maven-metadata.xml.sha1",
			expected: true,
		},
		{
			url:      "some/other/path/artifact.jar",
			expected: false,
		},
		{
			url:      "/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: false,
		},
	}

	for _, tc := range tests {
		got := IsMetadataURL(tc.url)
		if got != tc.expected {
			t.Errorf("IsMetadataURL(%q) = %v, expected %v", tc.url, got, tc.expected)
		}
	}
}

func TestFilenameFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "myapp-1.0.0.jar",
		},
		{
			url:      "/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "lang3-3.12.0.jar",
		},
	}

	for _, tc := range tests {
		filename, err := FilenameFromURL(tc.url)
		if err != nil {
			t.Errorf("Unexpected error for URL %s: %v", tc.url, err)
		}

		if filename != tc.expected {
			t.Errorf("For URL %s, expected filename %s but got %s", tc.url, tc.expected, filename)
		}
	}
}
