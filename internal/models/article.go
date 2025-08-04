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
	NumViews    int    `json:"views" db:"views"`
	NumLikes    int    `json:"num_likes" db:"num_likes"`
	NumComments int    `json:"num_comments" db:"num_comments"`
	IsLiked     bool   `json:"is_liked" db:"is_liked"`         // Whether the article is liked by the user
	IsAuthor    bool   `json:"is_author" db:"is_author"`       // Whether user sending request is the author
	IsFavorited bool   `json:"is_favorited" db:"is_favorited"` // Whether user sending request has added article to favorites
}

// DbArticle models how an article is stored in the database.
type DbArticle struct {
	ArticleID    string    `json:"article_id" db:"article_id"`
	Author       string    `json:"author_uid" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Content      string    `json:"content" db:"content"`
	Preview      string    `json:"preview" db:"preview"`
}

// Create a new article with a unique article ID
func (f *ArticleFactory) New(author string, title string, content string) *DbArticle {
	return &DbArticle{
		ArticleID:    "a" + gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		Author:       author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Flagged:      false,
		IsAvailable:  true,
		Content:      utils.SanitizeContent(content),
		Preview:      GetPreview(content),
	}
}
