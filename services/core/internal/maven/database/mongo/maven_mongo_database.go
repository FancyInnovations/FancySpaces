package mongo

import (
	"context"
	"errors"

	"github.com/fancyinnovations/fancyspaces/internal/maven"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DB struct {
	reposColl     *mongo.Collection
	artifactsColl *mongo.Collection
}

type Configuration struct {
	Mongo *mongo.Database
}

func NewDB(config *Configuration) *DB {
	reposColl := config.Mongo.Collection("maven_repositories")
	artifactsColl := config.Mongo.Collection("maven_artifacts")

	return &DB{
		reposColl:     reposColl,
		artifactsColl: artifactsColl,
	}
}

func (db *DB) GetRepository(ctx context.Context, spaceID, repoName string) (*maven.Repository, error) {
	filter := bson.D{{"space_id", spaceID}, {"name", repoName}}

	res := db.reposColl.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, maven.ErrRepositoryNotFound
		}
		return nil, res.Err()
	}

	var repo maven.Repository
	if err := res.Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (db *DB) GetRepositories(ctx context.Context, spaceID string) ([]maven.Repository, error) {
	filter := bson.D{{"space_id", spaceID}}

	cur, err := db.reposColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var repos []maven.Repository
	for cur.Next(ctx) {
		var repo maven.Repository
		if err := cur.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (db *DB) CreateRepository(ctx context.Context, repo maven.Repository) error {
	_, err := db.reposColl.InsertOne(ctx, repo)
	return err
}

func (db *DB) UpdateRepository(ctx context.Context, repo maven.Repository) error {
	filter := bson.D{{"space_id", repo.SpaceID}, {"name", repo.Name}}
	update := bson.D{{"$set", repo}}

	_, err := db.reposColl.UpdateOne(ctx, filter, update)
	return err
}

func (db *DB) DeleteRepository(ctx context.Context, spaceID, repoName string) error {
	filter := bson.D{{"space_id", spaceID}, {"name", repoName}}
	_, err := db.reposColl.DeleteOne(ctx, filter)
	return err
}

func (db *DB) GetArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) (*maven.Artifact, error) {
	filter := bson.D{{"space_id", spaceID}, {"repository", repoName}, {"group", groupID}, {"id", artifactID}}

	res := db.artifactsColl.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, maven.ErrArtifactNotFound
		}
		return nil, res.Err()
	}

	var artifact maven.Artifact
	if err := res.Decode(&artifact); err != nil {
		return nil, err
	}

	return &artifact, nil
}

func (db *DB) GetArtifacts(ctx context.Context, spaceID, repoName string) ([]maven.Artifact, error) {
	filter := bson.D{{"space_id", spaceID}, {"repository", repoName}}

	cur, err := db.artifactsColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var artifacts []maven.Artifact
	for cur.Next(ctx) {
		var artifact maven.Artifact
		if err := cur.Decode(&artifact); err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return artifacts, nil
}

func (db *DB) CreateArtifact(ctx context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	_, err := db.artifactsColl.InsertOne(ctx, artifact)
	return err
}

func (db *DB) UpdateArtifact(ctx context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	filter := bson.D{{"space_id", spaceID}, {"repository", repoName}, {"group", artifact.Group}, {"id", artifact.ID}}
	update := bson.D{{"$set", artifact}}

	_, err := db.artifactsColl.UpdateOne(ctx, filter, update)
	return err
}

func (db *DB) DeleteArtifact(ctx context.Context, spaceID, repoName, groupID, artifactID string) error {
	filter := bson.D{{"space_id", spaceID}, {"repository", repoName}, {"group", groupID}, {"id", artifactID}}
	_, err := db.artifactsColl.DeleteOne(ctx, filter)
	return err
}
