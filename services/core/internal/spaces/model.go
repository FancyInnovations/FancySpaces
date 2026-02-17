package spaces

import "github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"

type CreateOrUpdateSpaceReq struct {
	Slug        string            `json:"slug"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Categories  []spaces.Category `json:"categories"`
	IconURL     string            `json:"icon_url"`
}
