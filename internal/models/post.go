package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/microcosm-cc/bluemonday"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
)

type PostFactory struct {
	ContentFactory
}

func NewPostFactory() *PostFactory {
	return &PostFactory{}
}

type Post struct {
	PostID      string    `json:"post_id" db:"post_id"`
	Author      string    `json:"author" db:"author" binding:"required"`
	ThreadId    string    `json:"thread_id" db:"thread_id" binding:"required"`
	ReplyTo     *string   `json:"reply_to" db:"reply_to"`
	Title       string    `json:"title" db:"title" binding:"required"`
	PostContent string    `json:"content" db:"content" binding:"required"`
	TimeCreated time.Time `json:"time_created" db:"time_created"`
	LastEdited  time.Time `json:"last_edited" db:"last_edited"`
	Flagged     bool      `json:"flagged" db:"flagged"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	IsHeader    bool      `json:"is_header" db:"is_header" binding:"required"`

	// These columns are not in the database
	NumLikes int  `json:"num_likes" db:"-"`
	IsLiked  bool `json:"is_liked" db:"-"`
}

func (f *PostFactory) New(author string, threadId string, title string, content string, replyTo *string, isHeader bool) *Post {
	// Sanitize content to prevent XSS attacks
	policy := bluemonday.UGCPolicy()
	// Allow styles on images (to allow image resizing)
	policy.AllowStyles("width", "height", "draggable").OnElements("img")
	content = policy.Sanitize(content)

	return &Post{
		PostID:      "p" + gonanoid.Must(constants.CONTENT_ID_LENGTH),
		Author:      author,
		ThreadId:    threadId,
		ReplyTo:     replyTo,
		Title:       title,
		PostContent: content,
		TimeCreated: time.Now(),
		LastEdited:  time.Now(),
		Flagged:     false,
		IsAvailable: true,
		IsHeader:    isHeader,
	}
}

func (p *Post) GetID() string {
	return p.PostID
}

func (p *Post) GetAuthor() string {
	return p.Author
}

func (p *Post) GetTitle() string {
	return p.Title
}

func (p *Post) GetTimeCreated() time.Time {
	return p.TimeCreated
}

func (p *Post) GetLastActivity() time.Time {
	return p.LastEdited
}

func (p *Post) GetFlagged() bool {
	return p.Flagged
}
