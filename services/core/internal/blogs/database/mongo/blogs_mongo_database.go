package mongo

import (
	"context"
	"errors"

	"github.com/fancyinnovations/fancyspaces/core/internal/blogs"
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
	coll := config.Mongo.Collection("blogs")

	return &DB{
		coll: coll,
	}
}

func (db *DB) GetArticlesForSpace(spaceID string) ([]blogs.Article, error) {
	filter := bson.D{{"space_id", spaceID}}

	cur, err := db.coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var articles []blogs.Article
	for cur.Next(context.Background()) {
		var article blogs.Article
		if err := cur.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (db *DB) GetArticlesForUser(userID string) ([]blogs.Article, error) {
	filter := bson.D{{"author", userID}, {"space_id", ""}}

	cur, err := db.coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var articles []blogs.Article
	for cur.Next(context.Background()) {
		var article blogs.Article
		if err := cur.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (db *DB) GetArticleByID(articleID string) (*blogs.Article, error) {
	filter := bson.D{{"article_id", articleID}}

	res := db.coll.FindOne(context.Background(), filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, blogs.ErrArticleNotFound
		}
		return nil, res.Err()
	}

	var article blogs.Article
	if err := res.Decode(&article); err != nil {
		return nil, err
	}

	return &article, nil
}

func (db *DB) CreateArticle(article *blogs.Article) error {
	_, err := db.coll.InsertOne(context.Background(), article)
	return err
}

func (db *DB) UpdateArticle(article *blogs.Article) error {
	filter := bson.D{{"article_id", article.ID}}
	update := bson.D{{"$set", article}}

	res, err := db.coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return blogs.ErrArticleNotFound
	}

	return nil
}

func (db *DB) DeleteArticle(articleID string) error {
	filter := bson.D{{"article_id", articleID}}

	res, err := db.coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return blogs.ErrArticleNotFound
	}

	return nil
}
