package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ArticleFactory struct {
	ContentFactory
}

func NewArticleFactory() *ArticleFactory {
	return &ArticleFactory{}
}

// Article models how an article is retrieved from the database.
type Article struct {
	DbArticle

	Author      string `json:"author" db:"author_name"` // Name of the author
	NumLikes    int    `json:"num_likes" db:"num_likes"`
	IsLiked     bool   `json:"is_liked" db:"is_liked"` // Whether the article is liked by the user
	NumComments int    `json:"num_comments" db:"num_comments"`
}

// DbArticle models how an article is stored in the database.
type DbArticle struct {
	ArticleID    string    `json:"article_id" db:"article_id"`
	Author       string    `json:"author_uid" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Content      string    `json:"content" db:"content"`
}

// Create a new article with a unique article ID
func (f *ArticleFactory) New(author string, title string, content string) *DbArticle {
	return &DbArticle{
		ArticleID:    "a" + gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		Author:       author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Views:        0,
		Flagged:      false,
		IsAvailable:  true,
		Content:      utils.SanitizeContent(content),
	}
}

// Articles metadata
type ArticlesMetadata struct {
	NumArticles int `json:"num_articles" db:"num_articles"`
	NumPages    int `json:"num_pages" db:"-"`
}
