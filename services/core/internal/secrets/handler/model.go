package handler

type CreateSecretReq struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UpdateSecretReq struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}
