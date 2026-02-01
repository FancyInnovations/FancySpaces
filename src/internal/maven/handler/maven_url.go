package handler

import (
	"fmt"
	"strings"
)

func GroupFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 6 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	groupParts := parts[4 : len(parts)-3]

	return strings.Join(groupParts, "."), nil
}

func ArtifactFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 6 {
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

func FilenameFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 6 {
		return "", fmt.Errorf("invalid Maven URL: %s", url)
	}

	return parts[len(parts)-1], nil
}
