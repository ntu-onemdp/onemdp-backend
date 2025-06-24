package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
)

type ArticleFactory struct {
	ContentFactory
}

func NewArticleFactory() *ArticleFactory {
	return &ArticleFactory{}
}

type Article struct {
	ArticleID    string    `json:"article_id" db:"article_id"`
	Author       string    `json:"author" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Content      string    `json:"content" db:"content"`

	// Following fields are not stored in the database.
	NumLikes int `json:"num_likes" db:"-"`
}

// Create a new article with a unique article ID
func (f *ArticleFactory) New(author string, title string, content string) *Article {
	return &Article{
		ArticleID:    "a" + gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		Author:       author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Views:        0,
		Flagged:      false,
		IsAvailable:  true,
		Content:      content,
	}
}
