package handler

import "github.com/fancyinnovations/fancyspaces/storage/internal/database"

type CreateDatabaseReq struct {
	Name string `json:"name"`
}

type UpdateDatabaseReq struct {
	Users map[string]database.PermissionLevel `json:"users"`
}

type CreateCollectionReq struct {
	Name   string          `json:"name"`
	Engine database.Engine `json:"engine"`
}
