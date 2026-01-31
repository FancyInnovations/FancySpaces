package mongo

import (
	"errors"

	"github.com/fancyinnovations/fancyspaces/internal/issues"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DB struct {
	icoll *mongo.Collection
	ccoll *mongo.Collection
}

type Configuration struct {
	Mongo *mongo.Database
}

func NewDB(config *Configuration) *DB {
	icoll := config.Mongo.Collection("issues")
	ccoll := config.Mongo.Collection("comments")

	return &DB{
		icoll: icoll,
		ccoll: ccoll,
	}
}

func (db *DB) GetIssues(space string) ([]issues.Issue, error) {
	filter := bson.D{{"space", space}}

	cur, err := db.icoll.Find(nil, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(nil)

	var issuesList []issues.Issue
	for cur.Next(nil) {
		var issue issues.Issue
		if err := cur.Decode(&issue); err != nil {
			return nil, err
		}
		issuesList = append(issuesList, issue)
	}

	return issuesList, nil
}

func (db *DB) GetIssue(space, id string) (*issues.Issue, error) {
	filter := bson.D{
		{"space", space},
		{"id", id},
	}

	res := db.icoll.FindOne(nil, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, issues.ErrIssueNotFound
		}
		return nil, res.Err()
	}

	var issue issues.Issue
	if err := res.Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (db *DB) CreateIssue(issue *issues.Issue) error {
	_, err := db.icoll.InsertOne(nil, issue)
	return err
}

func (db *DB) UpdateIssue(issue *issues.Issue) error {
	filter := bson.D{
		{"space", issue.Space},
		{"id", issue.ID},
	}

	update := bson.D{
		{"$set", issue},
	}

	res, err := db.icoll.UpdateOne(nil, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return issues.ErrIssueNotFound
	}

	return nil
}

func (db *DB) DeleteIssue(space, id string) error {
	filter := bson.D{
		{"space", space},
		{"id", id},
	}

	_, err := db.icoll.DeleteOne(nil, filter)
	return err
}

func (db *DB) GetComments(issue string) ([]issues.Comment, error) {
	filter := bson.D{{"issue", issue}}

	cur, err := db.ccoll.Find(nil, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(nil)

	var comments []issues.Comment
	for cur.Next(nil) {
		var comment issues.Comment
		if err := cur.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *DB) AddComment(comment *issues.Comment) error {
	_, err := db.ccoll.InsertOne(nil, comment)
	return err
}

func (db *DB) UpdateComment(comment *issues.Comment) error {
	filter := bson.D{
		{"issue", comment.Issue},
		{"id", comment.ID},
	}

	update := bson.D{
		{"$set", comment},
	}

	res, err := db.ccoll.UpdateOne(nil, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return issues.ErrCommentNotFound
	}

	return nil
}

func (db *DB) DeleteComment(issue, id string) error {
	filter := bson.D{
		{"issue", issue},
		{"id", id},
	}

	_, err := db.ccoll.DeleteOne(nil, filter)
	return err
}
