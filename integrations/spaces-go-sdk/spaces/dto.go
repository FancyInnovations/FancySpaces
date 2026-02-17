package spaces

type GetSecretReqDTO struct {
	SpaceID string `json:"space_id"`
	Key     string `json:"key"`
}

type GetSecretRespDTO struct {
	Value string `json:"value"`
}
