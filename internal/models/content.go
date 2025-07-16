package models

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

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

// Utility function to get preview from content
func GetPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 100

	p := bluemonday.StrictPolicy()
	content = p.Sanitize(content)

	utils.Logger.Debug().Str("content", content).Msg("Sanitized content")
	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}
