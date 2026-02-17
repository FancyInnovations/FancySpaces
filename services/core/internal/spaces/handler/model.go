package handler

import (
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

type ChangeStatusReq struct {
	To spaces.Status `json:"to"`
}

type SpaceDownloadsResp struct {
	Downloads uint64            `json:"downloads"`
	Versions  map[string]uint64 `json:"versions"`
}
