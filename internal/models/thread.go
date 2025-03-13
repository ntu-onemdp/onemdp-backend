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

type Thread struct {
	Content
	ThreadID     string    `json:"thread_id" db:"thread_id"`
	Author       string    `json:"author" db:"author"`
	Title        string    `json:"title" db:"title"`
	TimeCreated  time.Time `json:"time_created" db:"time_created"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Views        int       `json:"views" db:"views"`
	Flagged      bool      `json:"flagged" db:"flagged"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	Preview      string    `json:"preview" db:"preview"`
}

// Create a new thread with a unique thread ID
func (f *ThreadFactory) New(author string, title string, content string) *Thread {
	return &Thread{
		ThreadID:     gonanoid.Must(constants.CONTENT_ID_LENGTH), // Note that this can cause program to panic!
		Author:       author,
		Title:        title,
		TimeCreated:  time.Now(),
		LastActivity: time.Now(),
		Views:        0,
		Flagged:      false,
		IsAvailable:  true,
		Preview:      getPreview(content),
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
func getPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 100

	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}
