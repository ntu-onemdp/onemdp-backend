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

type ContentMetadata struct {
	Total    int `json:"total" db:"count"`
	NumPages int `json:"num_pages" db:"-"`
}

// Use this interface to create new content objects
type ContentFactory interface {
	New(any) *Content
}

// Utility function to get preview from content
func GetPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 250

	p := bluemonday.StrictPolicy()
	content = p.Sanitize(content)

	utils.Logger.Trace().Str("content", content).Msg("Sanitized content")
	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}

// Column definitions available for sorting
type SortColumn string

const (
	TIME_CREATED_COL  SortColumn = "time_created"
	LAST_ACTIVITY_COL SortColumn = "last_activity"
	VIEWS_COL         SortColumn = "views"
)

// Convert string to SortColumn
func StrToSortColumn(s string) SortColumn {
	switch s {
	case "time_created":
		return TIME_CREATED_COL
	case "last_activity":
		return LAST_ACTIVITY_COL
	case "views":
		return VIEWS_COL
	default:
		return TIME_CREATED_COL
	}
}
