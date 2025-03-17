package models

import "time"

// The content interface defines the methods that must be implemented by the following:
// - Article
// - Comment
// - Thread
// - Post
type Content interface {
	GetID() string
	GetAuthor() string
	GetTitle() string // Returns content for comment
	GetTimeCreated() time.Time
	GetLastActivity() time.Time
	GetFlagged() bool
}

// Use this interface to create new content objects
type ContentFactory interface {
	New(any) *Content
}
