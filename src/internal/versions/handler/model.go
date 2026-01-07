package handler

import "github.com/fancyinnovations/fancyspaces/internal/versions"

type CreateVersionReq struct {
	Name                      string            `json:"name"`
	Platform                  versions.Platform `json:"platform"`
	Channel                   versions.Channel  `json:"channel"`
	Changelog                 string            `json:"changelog"`
	SupportedPlatformVersions []string          `json:"supported_platform_versions"`
}

type VersionDownloadsResp struct {
	Downloads uint64 `json:"downloads"`
}
