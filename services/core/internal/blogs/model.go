package blogs

import "time"

type Article struct {
	ID      string `json:"id" bson:"article_id"`
	SpaceID string `json:"space_id" bson:"space_id"` // if space_id is empty, this article is user-owned, otherwise it's space-owned
	Author  string `json:"author" bson:"author"`

	Title   string `json:"title" bson:"title"`
	Summary string `json:"summary" bson:"summary"`
	Content string `json:"-" bson:"content"`

	PublishedAt time.Time `json:"published_at" bson:"published_at"`
}
