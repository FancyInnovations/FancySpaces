package handler

import (
	"fmt"
	"strings"
	"unicode"
)

// looksLikeVersion returns true if the path segment appears to be a version.
// Heuristic: contains a digit or the word SNAPSHOT.
func looksLikeVersion(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	if strings.Contains(s, "SNAPSHOT") {
		return true
	}
	return false
}

func GroupFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")

	if IsMetadataURL(url) {
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
		}

		// If the segment before the metadata filename looks like a version, the layout
		// is .../{groupPath}/{artifactId}/{version}/maven-metadata.xml
		if looksLikeVersion(parts[len(parts)-2]) {
			if len(parts) < 6 {
				return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
			}
			groupParts := parts[3 : len(parts)-3]
			return strings.Join(groupParts, "."), nil
		}

		// artifact-level metadata: .../{groupPath}/{artifactId}/maven-metadata.xml
		groupParts := parts[3 : len(parts)-2]
		return strings.Join(groupParts, "."), nil
	}

	// normal artifact file URL: .../{groupPath}/{artifactId}/{version}/{filename}
	if len(parts) < 6 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	groupParts := parts[3 : len(parts)-3]
	return strings.Join(groupParts, "."), nil
}

func ArtifactFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")

	if IsMetadataURL(url) {
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid Maven metadata URL: %s", url)
		}

		// if the segment before the metadata filename looks like a version,
		// the artifact id is one segment earlier
		if looksLikeVersion(parts[len(parts)-2]) {
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
	parts := strings.Split(url, "/")
	if len(parts) < 6 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-2], nil
}

func IsMetadataURL(url string) bool {
	return strings.HasSuffix(url, "maven-metadata.xml") || strings.Contains(url, "maven-metadata.xml.")
}

func FilenameFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-1], nil
}
