package models

import "time"

type Like struct {
	Username    string    `json:"username" db:"username"`
	ContentId   string    `json:"content_id" db:"content_id"`     // Parse to UUID if content type is post or comment
	ContentType string    `json:"content_type" db:"content_type"` // thread, article, post, comment
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

// Content types
const THREAD_CONTENT_TYPE = "thread"
const ARTICLE_CONTENT_TYPE = "article"
const POST_CONTENT_TYPE = "post"
const COMMENT_CONTENT_TYPE = "comment"
