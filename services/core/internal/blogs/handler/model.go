package handler

type CreateOrUpdateArticleReq struct {
	SpaceID string `json:"space_id,omitempty"`

	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}
