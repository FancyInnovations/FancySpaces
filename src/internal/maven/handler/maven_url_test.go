package handler

import "testing"

func TestGroupFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/maven/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "com.example.myapp",
		},
		{
			url:      "/maven/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "org.apache.commons.lang3",
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
			url:      "/maven/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "myapp",
		},
		{
			url:      "/maven/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
			expected: "lang3",
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
			url:      "/maven/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "1.0.0",
		},
		{
			url:      "/maven/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
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

func TestFilenameFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "/maven/space1/repo1/com/example/myapp/1.0.0/myapp-1.0.0.jar",
			expected: "myapp-1.0.0.jar",
		},
		{
			url:      "/maven/space2/repo2/org/apache/commons/lang3/3.12.0/lang3-3.12.0.jar",
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
