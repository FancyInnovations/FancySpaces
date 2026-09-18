package mongo

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	ctx := context.Background()
	_, err := icoll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "space", Value: 1}, {Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "space", Value: 1}, {Key: "status", Value: 1}, {Key: "updated_at", Value: -1}}},
		{Keys: bson.D{{Key: "space", Value: 1}, {Key: "external_source", Value: 1}, {Key: "external_id", Value: 1}}},
	})
	if err != nil {
		slog.Error("failed to create issue indexes", "error", err)
	}
	_, err = ccoll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "space", Value: 1}, {Key: "issue", Value: 1}, {Key: "created_at", Value: 1}}},
		{Keys: bson.D{{Key: "space", Value: 1}, {Key: "issue", Value: 1}, {Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		slog.Error("failed to create comment indexes", "error", err)
	}

	return &DB{
		icoll: icoll,
		ccoll: ccoll,
	}
}

func (db *DB) GetIssues(space string) ([]issues.Issue, error) {
	filter := bson.D{{"space", space}}

	ctx := context.Background()
	cur, err := db.icoll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var issuesList []issues.Issue
	for cur.Next(ctx) {
		var issue issues.Issue
		if err := cur.Decode(&issue); err != nil {
			return nil, err
		}
		issuesList = append(issuesList, issue)
	}

	return issuesList, nil
}

func (db *DB) ListIssues(space string, opts issues.ListOptions) ([]issues.Issue, int, error) {
	filter := bson.D{{"space", space}, {"archived_at", nil}}
	if opts.Type != "" {
		filter = append(filter, bson.E{Key: "type", Value: opts.Type})
	}
	if opts.Status != "" {
		filter = append(filter, bson.E{Key: "status", Value: opts.Status})
	}
	if opts.Priority != "" {
		filter = append(filter, bson.E{Key: "priority", Value: opts.Priority})
	}
	if opts.Assignee != "" {
		filter = append(filter, bson.E{Key: "assignee", Value: opts.Assignee})
	}
	if opts.ExternalSource != "" {
		filter = append(filter, bson.E{Key: "external_source", Value: opts.ExternalSource})
	}
	if query := strings.TrimSpace(opts.Query); query != "" {
		pattern := bson.Regex{Pattern: regexp.QuoteMeta(query), Options: "i"}
		filter = append(filter, bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "id", Value: pattern}},
			bson.D{{Key: "title", Value: pattern}},
			bson.D{{Key: "description", Value: pattern}},
		}})
	}

	ctx := context.Background()
	total64, err := db.icoll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	findOpts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}).SetSkip(int64(opts.Offset)).SetLimit(int64(opts.Limit))
	cur, err := db.icoll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	result := make([]issues.Issue, 0)
	for cur.Next(ctx) {
		var issue issues.Issue
		if err := cur.Decode(&issue); err != nil {
			return nil, 0, err
		}
		result = append(result, issue)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return result, int(total64), nil
}

func (db *DB) GetIssue(space, id string) (*issues.Issue, error) {
	filter := bson.D{
		{"space", space},
		{"id", id},
		{"archived_at", nil},
	}

	res := db.icoll.FindOne(context.Background(), filter)
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
	_, err := db.icoll.InsertOne(context.Background(), issue)
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

	res, err := db.icoll.UpdateOne(context.Background(), filter, update)
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

	_, err := db.icoll.DeleteOne(context.Background(), filter)
	return err
}

func (db *DB) GetComments(space, issue string) ([]issues.Comment, error) {
	filter := bson.D{{"space", space}, {"issue", issue}}

	cur, err := db.ccoll.Find(context.Background(), filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	defer cur.Close(ctx)

	var comments []issues.Comment
	for cur.Next(ctx) {
		var comment issues.Comment
		if err := cur.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *DB) AddComment(comment *issues.Comment) error {
	_, err := db.ccoll.InsertOne(context.Background(), comment)
	return err
}

func (db *DB) UpdateComment(comment *issues.Comment) error {
	filter := bson.D{
		{"space", comment.Space},
		{"issue", comment.Issue},
		{"id", comment.ID},
	}

	update := bson.D{
		{"$set", comment},
	}

	res, err := db.ccoll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return issues.ErrCommentNotFound
	}

	return nil
}

func (db *DB) DeleteComment(space, issue, id string) error {
	filter := bson.D{
		{"space", space},
		{"issue", issue},
		{"id", id},
	}

	_, err := db.ccoll.DeleteOne(context.Background(), filter)
	return err
}
