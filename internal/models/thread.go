package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
)

type ThreadFactory struct {
	ContentFactory
}

func NewThreadFactory() *ThreadFactory {
	return &ThreadFactory{}
}

// Thread models how a thread is retrieved from the database.
type Thread struct {
	DbThread

	Author     string `json:"author" db:"author_name"` // Name of the author
	NumLikes   int    `json:"num_likes" db:"num_likes"`
	NumReplies int    `json:"num_replies" db:"num_replies"`
	IsLiked    bool   `json:"is_liked" db:"is_liked"`   // Whether the thread is liked by the user
	IsAuthor   bool   `json:"is_author" db:"is_author"` // Whether user sending request is the author
}

// DbThread models how a thread is stored in the database.
type DbThread struct {
	ThreadID     string    `json:"thread_id" db:"thread_id"`
	AuthorUid    string    `json:"author_uid" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Preview      string    `json:"preview" db:"preview"`
	IsAnon       bool      `json:"is_anon" db:"is_anon"`
}

// Create a new thread with a unique thread ID
func (f *ThreadFactory) New(author string, title string, content string, isAnon bool) *DbThread {
	return &DbThread{
		ThreadID:     "t" + gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		AuthorUid:    author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Views:        0,
		Flagged:      false,
		IsAvailable:  true,
		Preview:      GetPreview(content),
		IsAnon:       isAnon,
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
