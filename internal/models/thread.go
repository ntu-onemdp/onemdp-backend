package models

import (
	"time"
)

type Thread struct {
	ThreadId     string    `json:"thread_id" db:"thread_id"`
	Author       string    `json:"author" db:"author"`
	Title        string    `json:"title" db:"title"`
	NumLikes     int       `json:"num_likes" db:"num_likes"`
	NumReplies   int       `json:"num_replies" db:"num_replies"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Preview      string    `json:"preview" db:"preview"`
}

// NewThread has minimal fields. Database takes care of the default field values.
type NewThread struct {
	Author  string `json:"author" db:"author"`
	Title   string `json:"title" db:"title"`
	Preview string `json:"preview" db:"preview"`
}
