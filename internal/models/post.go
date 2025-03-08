package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	PostId      uuid.UUID `json:"post_id" db:"post_id"`
	Author      string    `json:"author" db:"author"`
	ThreadId    uuid.UUID `json:"thread_id" db:"thread_id"`
	ReplyTo     *string   `json:"reply_to" db:"reply_to"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	NumLikes    int       `json:"num_likes" db:"num_likes"`
	TimeCreated time.Time `json:"time_created" db:"time_created"`
	LastEdited  time.Time `json:"last_edited" db:"last_edited"`
	Flagged     bool      `json:"flagged" db:"flagged"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
}

// NewPost has minimal fields. Database takes care of the default field values.
type NewPost struct {
	Author   string    `json:"author" db:"author"`
	ThreadId uuid.UUID `json:"thread_id" db:"thread_id"`
	Title    string    `json:"title" db:"title"`
	Content  string    `json:"content" db:"content"`
	ReplyTo  *string   `json:"reply_to" db:"reply_to"`
}
