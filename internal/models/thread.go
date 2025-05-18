package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/microcosm-cc/bluemonday"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadFactory struct {
	ContentFactory
}

func NewThreadFactory() *ThreadFactory {
	return &ThreadFactory{}
}

type Thread struct {
	ThreadID     string    `json:"thread_id" db:"thread_id"`
	Author       string    `json:"author" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Preview      string    `json:"preview" db:"preview"`

	// Following fields are not stored in the database
	NumLikes   int  `json:"num_likes" db:"-"`
	NumReplies int  `json:"num_replies" db:"-"`
	IsLiked    bool `json:"is_liked" db:"-"` // Whether the thread is liked by the user
}

// Create a new thread with a unique thread ID
func (f *ThreadFactory) New(author string, title string, content string) *Thread {
	return &Thread{
		ThreadID:     "t" + gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		Author:       author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Views:        0,
		Flagged:      false,
		IsAvailable:  true,
		Preview:      GetPreview(content),
	}
}

// Column definitions available for sorting
type ThreadColumn string

const (
	TIME_CREATED_COL  ThreadColumn = "time_created"
	LAST_ACTIVITY_COL ThreadColumn = "last_activity"
)

// Threads metadata
type ThreadsMetadata struct {
	NumThreads int `json:"num_threads"`
}

// Convert string to ThreadColumn
func StrToThreadColumn(s string) ThreadColumn {
	switch s {
	case "time_created":
		return TIME_CREATED_COL
	case "last_activity":
		return LAST_ACTIVITY_COL
	default:
		return TIME_CREATED_COL
	}
}

func (t *Thread) GetID() string {
	return t.ThreadID
}

func (t *Thread) GetAuthor() string {
	return t.Author
}

func (t *Thread) GetTitle() string {
	return t.Title
}

func (t *Thread) GetTimeCreated() time.Time {
	return t.TimeCreated
}

func (t *Thread) GetLastActivity() time.Time {
	return t.LastActivity
}

func (t *Thread) GetFlagged() bool {
	return t.Flagged
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
