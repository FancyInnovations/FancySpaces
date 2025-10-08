package handler

import "github.com/fancyinnovations/fancyspaces/src/internal/spaces"

type ChangeStatusReq struct {
	To spaces.Status `json:"to"`
}
