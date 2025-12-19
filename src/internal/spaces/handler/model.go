package handler

import "github.com/fancyinnovations/fancyspaces/internal/spaces"

type ChangeStatusReq struct {
	To spaces.Status `json:"to"`
}
