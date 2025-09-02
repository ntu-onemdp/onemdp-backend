package models

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
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

	// Define line break tags to look for (lowercase for case-insensitive)
	lineBreakTags := []string{"</p>", "</h1>", "</h2>", "</h3>", "</h4>", "</h5>", "</h6>", "<br>", "<br/>", "<br />"}

	lowerContent := strings.ToLower(content)

	// Find the earliest occurrence of any line break tag
	earliestIdx := -1
	for _, tag := range lineBreakTags {
		idx := strings.Index(lowerContent, tag)
		if idx != -1 && (earliestIdx == -1 || idx < earliestIdx) {
			earliestIdx = idx + len(tag) // cut after the tag
		}
	}

	if earliestIdx != -1 {
		// Cut content abruptly at earliest line break tag
		cutContent := content[:earliestIdx]

		// Sanitize truncated content
		p := bluemonday.StrictPolicy()
		sanitized := p.Sanitize(cutContent)
		return strings.TrimSpace(sanitized)
	}

	// No line break found: sanitize full content and truncate nicely
	p := bluemonday.StrictPolicy()
	sanitized := p.Sanitize(content)
	sanitized = strings.TrimSpace(sanitized)

	if utf8.RuneCountInString(sanitized) <= MAX_PREVIEW_LENGTH {
		return sanitized
	}

	truncated := []rune(sanitized)[:MAX_PREVIEW_LENGTH]
	preview := string(truncated)

	lastSpace := strings.LastIndex(preview, " ")
	if lastSpace > 0 {
		preview = preview[:lastSpace]
	}

	return preview + "…"
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
