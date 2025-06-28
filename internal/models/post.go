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

// Post models how a post is retrieved from the database.
type Post struct {
	DbPost

	Author   string `json:"author" db:"author_name"` // Name of the author
	NumLikes int    `json:"num_likes" db:"num_likes"`
	IsLiked  bool   `json:"is_liked" db:"is_liked"` // Whether the post is liked by the user
}

// DbPost models how a post is stored in the database.
type DbPost struct {
	PostID      string    `json:"post_id" db:"post_id"`
	AuthorUid   string    `json:"author_uid" db:"author" binding:"required"`
	ThreadId    string    `json:"thread_id" db:"thread_id" binding:"required"`
	ReplyTo     *string   `json:"reply_to" db:"reply_to"`
	Title       string    `json:"title" db:"title" binding:"required"`
	PostContent string    `json:"content" db:"content" binding:"required"`
	TimeCreated time.Time `json:"time_created" db:"time_created"`
	LastEdited  time.Time `json:"last_edited" db:"last_edited"`
	Flagged     bool      `json:"flagged" db:"flagged"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	IsHeader    bool      `json:"is_header" db:"is_header" binding:"required"`
}

func (f *PostFactory) New(author string, threadId string, title string, content string, replyTo *string, isHeader bool) *DbPost {
	// Sanitize content to prevent XSS attacks
	policy := bluemonday.UGCPolicy()
	// Allow styles on images (to allow image resizing)
	policy.AllowStyles("width", "height", "draggable").OnElements("img")
	content = policy.Sanitize(content)

	return &DbPost{
		PostID:      "p" + gonanoid.Must(constants.CONTENT_ID_LENGTH),
		AuthorUid:   author,
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

func (p *DbPost) GetID() string {
	return p.PostID
}

func (p *DbPost) GetAuthor() string {
	return p.AuthorUid
}

func (p *DbPost) GetTitle() string {
	return p.Title
}

func (p *DbPost) GetTimeCreated() time.Time {
	return p.TimeCreated
}

func (p *DbPost) GetLastActivity() time.Time {
	return p.LastEdited
}

func (p *DbPost) GetFlagged() bool {
	return p.Flagged
}
