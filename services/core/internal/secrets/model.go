package secrets

import "time"

type Secret struct {
	SpaceID           string `json:"space_id" bson:"space_id"`
	Key               string `json:"key" bson:"key"`
	Value             []byte `json:"-" bson:"value"`
	DataEncryptionKey []byte `json:"-" bson:"data_encryption_key"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	Description string `json:"description" bson:"description"`
}
