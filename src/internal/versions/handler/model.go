package handler

type CreateVersionReq struct {
	SpaceID                   string   `json:"space_id"`
	Name                      string   `json:"name"`
	Channel                   string   `json:"channel"`
	Changelog                 string   `json:"changelog"`
	SupportedPlatformVersions []string `json:"supported_platform_versions"`
}
