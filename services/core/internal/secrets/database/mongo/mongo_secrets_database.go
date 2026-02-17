package mongo

import (
	"context"

	"github.com/fancyinnovations/fancyspaces/core/internal/secrets"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DB struct {
	coll *mongo.Collection
}

type Configuration struct {
	Mongo *mongo.Database
}

func NewDB(config *Configuration) *DB {
	coll := config.Mongo.Collection("secrets")

	return &DB{
		coll: coll,
	}
}

func (db *DB) GetSecret(spaceID, key string) (*secrets.Secret, error) {
	filter := bson.D{{"space_id", spaceID}, {"key", key}}

	res := db.coll.FindOne(context.Background(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, secrets.ErrSecretNotFound
		}
		return nil, res.Err()
	}

	var secret secrets.Secret
	if err := res.Decode(&secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

func (db *DB) GetSecrets(spaceID string) ([]*secrets.Secret, error) {
	filter := bson.D{{"space_id", spaceID}}

	cur, err := db.coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var secretsList []*secrets.Secret
	for cur.Next(context.Background()) {
		var secret secrets.Secret
		if err := cur.Decode(&secret); err != nil {
			return nil, err
		}
		secretsList = append(secretsList, &secret)
	}

	return secretsList, nil
}

func (db *DB) CreateSecret(secret *secrets.Secret) error {
	_, err := db.coll.InsertOne(context.Background(), secret)
	return err
}

func (db *DB) UpdateSecret(secret *secrets.Secret) error {
	filter := bson.D{{"space_id", secret.SpaceID}, {"key", secret.Key}}
	update := bson.D{{"$set", secret}}

	res, err := db.coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return secrets.ErrSecretNotFound
	}

	return nil
}

func (db *DB) DeleteSecret(spaceID, key string) error {
	filter := bson.D{{"space_id", spaceID}, {"key", key}}

	_, err := db.coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
