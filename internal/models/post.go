package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type PostFactory struct {
	ContentFactory
}

func NewPostFactory() *PostFactory {
	return &PostFactory{}
}

// Post models how a post is retrieved from the database.
type Post struct {
	DbPost

	Author   string `json:"author" db:"author_name"` // Name of the author
	NumLikes int    `json:"num_likes" db:"num_likes"`
	IsLiked  bool   `json:"is_liked" db:"is_liked"`   // Whether the post is liked by the user
	IsAuthor bool   `json:"is_author" db:"is_author"` // Whether user sending request is the author
}

// DbPost models how a post is stored in the database.
type DbPost struct {
	PostID           string           `json:"post_id" db:"post_id"`
	AuthorUid        string           `json:"author_uid" db:"author" binding:"required"`
	ThreadId         string           `json:"thread_id" db:"thread_id" binding:"required"`
	ReplyTo          *string          `json:"reply_to" db:"reply_to"`
	Title            string           `json:"title" db:"title" binding:"required"`
	PostContent      string           `json:"content" db:"content" binding:"required"`
	TimeCreated      time.Time        `json:"time_created" db:"time_created"`
	LastEdited       time.Time        `json:"last_edited" db:"last_edited"`
	Flagged          bool             `json:"flagged" db:"flagged"`
	IsAvailable      bool             `json:"is_available" db:"is_available"`
	IsHeader         bool             `json:"is_header" db:"is_header"`
	IsAnon           bool             `json:"is_anon" db:"is_anon"`
	ValidationStatus ValidationStatus `json:"validation_status" db:"validation_status"`
	ValidatedBy      *string          `json:"validated_by" db:"validated_by"`
}

func (f *PostFactory) New(author string, threadId string, title string, content string, replyTo *string, isHeader bool, isAnon bool) *DbPost {
	return &DbPost{
		PostID:           "p" + gonanoid.Must(constants.CONTENT_ID_LENGTH),
		AuthorUid:        author,
		ThreadId:         threadId,
		ReplyTo:          replyTo,
		Title:            title,
		PostContent:      utils.SanitizeContent(content),
		TimeCreated:      time.Now(),
		LastEdited:       time.Now(),
		Flagged:          false,
		IsAvailable:      true,
		IsHeader:         isHeader,
		IsAnon:           isAnon,
		ValidationStatus: ValidationUnverified,
	}
}

// Validation status (migration 20250916073659_posts_add_validation_status)
type ValidationStatus string

const (
	ValidationUnverified ValidationStatus = "unverified"
	ValidationValidated  ValidationStatus = "validated"
	ValidationRefuted    ValidationStatus = "refuted"
)

func (vs ValidationStatus) String() string {
	return string(vs)
}

func ParseValidationStatus(s string) (ValidationStatus, bool) {
	switch s {
	case "unverified":
		return ValidationUnverified, true
	case "validated":
		return ValidationValidated, true
	case "refuted":
		return ValidationRefuted, true
	default:
		return "", false
	}
}
