package blogs

import (
	"time"

	"github.com/OliverSchlueter/goutils/idgen"
)

type Database interface {
	GetArticlesForSpace(spaceID string) ([]Article, error)
	GetArticlesForUser(userID string) ([]Article, error)

	GetArticleByID(articleID string) (*Article, error)

	CreateArticle(article *Article) error
	UpdateArticle(article *Article) error
	DeleteArticle(articleID string) error
}

type Store struct {
	db Database
}

type Configuration struct {
	Database Database
}

func NewStore(config *Configuration) *Store {
	return &Store{
		db: config.Database,
	}
}

func (s *Store) GetArticlesForSpace(spaceID string) ([]Article, error) {
	return s.db.GetArticlesForSpace(spaceID)
}

func (s *Store) GetArticlesForUser(userID string) ([]Article, error) {
	return s.db.GetArticlesForUser(userID)
}

func (s *Store) GetArticleByID(articleID string) (*Article, error) {
	return s.db.GetArticleByID(articleID)
}

func (s *Store) GetArticleContentByID(articleID string) (string, error) {
	article, err := s.db.GetArticleByID(articleID)
	if err != nil {
		return "", err
	}

	return article.Content, nil
}

func (s *Store) CreateArticle(article *Article, content string) error {
	if len(article.Title) > MaxTitleSize {
		return ErrTitleTooLong
	}
	if len(article.Summary) > MaxSummarySize {
		return ErrSummaryTooLong
	}
	if len(content) > MaxContentSize {
		return ErrContentTooLong
	}

	article.ID = idgen.GenerateID(12)
	article.PublishedAt = time.Now()
	article.Content = content

	return s.db.CreateArticle(article)
}

func (s *Store) UpdateArticle(article *Article) error {
	if len(article.Title) > MaxTitleSize {
		return ErrTitleTooLong
	}
	if len(article.Summary) > MaxSummarySize {
		return ErrSummaryTooLong
	}
	if len(article.Content) > MaxContentSize {
		return ErrContentTooLong
	}

	return s.db.UpdateArticle(article)
}

func (s *Store) UpdateArticleContent(articleID string, content string) error {
	if len(content) > MaxContentSize {
		return ErrContentTooLong
	}

	article, err := s.db.GetArticleByID(articleID)
	if err != nil {
		return err
	}

	article.Content = content
	return s.db.UpdateArticle(article)
}

func (s *Store) DeleteArticle(articleID string) error {
	return s.db.DeleteArticle(articleID)
}
