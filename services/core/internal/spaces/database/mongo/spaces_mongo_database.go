package mongo

import (
	"context"
	"errors"

	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
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
	coll := config.Mongo.Collection("spaces")

	return &DB{
		coll: coll,
	}
}

func (db *DB) GetByID(id string) (*spaces.Space, error) {
	filter := bson.D{{"space_id", id}}

	res := db.coll.FindOne(context.Background(), filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, spaces.ErrSpaceNotFound
		}
		return nil, res.Err()
	}

	var sp spaces.Space
	if err := res.Decode(&sp); err != nil {
		return nil, err
	}

	return &sp, nil
}

func (db *DB) GetBySlug(slug string) (*spaces.Space, error) {
	filter := bson.D{{"slug", slug}}

	res := db.coll.FindOne(context.Background(), filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, spaces.ErrSpaceNotFound
		}
		return nil, res.Err()
	}

	var sp spaces.Space
	if err := res.Decode(&sp); err != nil {
		return nil, err
	}

	return &sp, nil
}

func (db *DB) GetAll() ([]spaces.Space, error) {
	cur, err := db.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var spacesList []spaces.Space
	for cur.Next(context.Background()) {
		var sp spaces.Space
		if err := cur.Decode(&sp); err != nil {
			return nil, err
		}
		spacesList = append(spacesList, sp)
	}

	return spacesList, nil
}

func (db *DB) Create(s *spaces.Space) error {
	_, err := db.coll.InsertOne(context.Background(), s)
	return err
}

func (db *DB) Update(id string, s *spaces.Space) error {
	filter := bson.D{{"space_id", id}}
	update := bson.D{{"$set", s}}

	res, err := db.coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return spaces.ErrSpaceNotFound
	}

	return nil
}

func (db *DB) Delete(id string) error {
	filter := bson.D{{"space_id", id}}
	_, err := db.coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
