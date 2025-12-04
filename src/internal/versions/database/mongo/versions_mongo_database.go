package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/fancyinnovations/fancyspaces/src/internal/versions"
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
	coll := config.Mongo.Collection("versions")

	return &DB{
		coll: coll,
	}
}

func (db *DB) GetAll(ctx context.Context, spaceID string) ([]versions.Version, error) {
	filter := bson.D{{"space_id", spaceID}}

	cur, err := db.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not find versions: %w", err)
	}
	defer cur.Close(ctx)

	var vers []versions.Version
	for cur.Next(ctx) {
		var v versions.Version
		if err := cur.Decode(&v); err != nil {
			return nil, fmt.Errorf("could not decode version: %w", err)
		}
		vers = append(vers, v)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return vers, nil
}

func (db *DB) GetByID(ctx context.Context, spaceID, versionID string) (*versions.Version, error) {
	filter := bson.D{{"space_id", spaceID}, {"id", versionID}}

	res := db.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, versions.ErrVersionNotFound
		}
		return nil, fmt.Errorf("could not find version by id: %w", res.Err())
	}

	var v versions.Version
	if err := res.Decode(&v); err != nil {
		return nil, fmt.Errorf("could not decode version: %w", err)
	}

	return &v, nil
}

func (db *DB) GetByName(ctx context.Context, spaceID, versionNumber string) (*versions.Version, error) {
	filter := bson.D{{"space_id", spaceID}, {"name", versionNumber}}

	res := db.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, versions.ErrVersionNotFound
		}
		return nil, fmt.Errorf("could not find version by name: %w", res.Err())
	}

	var v versions.Version
	if err := res.Decode(&v); err != nil {
		return nil, fmt.Errorf("could not decode version: %w", err)
	}

	return &v, nil
}

func (db *DB) Create(ctx context.Context, v *versions.Version) error {
	_, err := db.coll.InsertOne(ctx, v)
	if err != nil {
		return fmt.Errorf("could not insert version: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, spaceID, versionID string, v *versions.Version) error {
	filter := bson.D{{"space_id", spaceID}, {"id", versionID}}
	update := bson.D{{"$set", v}}

	res, err := db.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("could not update version: %w", err)
	}
	if res.MatchedCount == 0 {
		return versions.ErrVersionNotFound
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, spaceID, versionID string) error {
	filter := bson.D{{"space_id", spaceID}, {"id", versionID}}

	res, err := db.coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete version: %w", err)
	}
	if res.DeletedCount == 0 {
		return versions.ErrVersionNotFound
	}

	return nil
}

func (db *DB) LogDownload(ctx context.Context, spaceID, versionID string) error {
	filter := bson.D{{"space_id", spaceID}, {"id", versionID}}
	update := bson.D{{"$inc", bson.D{{"downloads", 1}}}}

	res, err := db.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("could not log download: %w", err)
	}
	if res.MatchedCount == 0 {
		return versions.ErrVersionNotFound
	}

	return nil
}
