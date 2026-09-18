package handler

import (
	"fmt"
	"strings"
)

// IsMetadataURL checks if a URL is for maven metadata (or its checksums).
func IsMetadataURL(url string) bool {
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return false
	}
	last := parts[len(parts)-1]
	return last == "maven-metadata.xml" || strings.HasPrefix(last, "maven-metadata.xml.")
}

func isSnapshotVersion(s string) bool {
	return strings.HasSuffix(s, "-SNAPSHOT")
}

func splitURL(url string) []string {
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	parts := strings.Split(url, "/")
	for len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	return parts
}

func GroupFromURL(url string) (string, error) {
	parts := splitURL(url)

	if IsMetadataURL(url) {
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
		}
		if isSnapshotVersion(parts[len(parts)-2]) {
			// version-level metadata: [space, repo, group..., artifact, version, metadata]
			if len(parts) < 5 {
				return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
			}
			return strings.Join(parts[2:len(parts)-3], "."), nil
		}
		// artifact-level metadata: [space, repo, group..., artifact, metadata]
		return strings.Join(parts[2:len(parts)-2], "."), nil
	}

	if len(parts) < 5 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return strings.Join(parts[2:len(parts)-3], "."), nil
}

func ArtifactFromURL(url string) (string, error) {
	parts := splitURL(url)

	if IsMetadataURL(url) {
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
		}
		if isSnapshotVersion(parts[len(parts)-2]) {
			if len(parts) < 5 {
				return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
			}
			return parts[len(parts)-3], nil
		}
		return parts[len(parts)-2], nil
	}

	if len(parts) < 5 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-3], nil
}

func VersionFromURL(url string) (string, error) {
	parts := splitURL(url)

	if IsMetadataURL(url) {
		if len(parts) >= 5 && isSnapshotVersion(parts[len(parts)-2]) {
			return parts[len(parts)-2], nil
		}
		return "", fmt.Errorf("metadata URL has no version: %s", url)
	}

	if len(parts) < 5 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-2], nil
}

func FilenameFromURL(url string) (string, error) {
	parts := splitURL(url)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-1], nil
}
