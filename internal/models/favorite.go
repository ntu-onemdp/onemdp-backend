package models

import "time"

type Favorite struct {
	Uid       string    `json:"uid" db:"uid"`
	ContentId string    `json:"content_id" db:"content_id"` // Parse to UUID if content type is post or comment
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewFavorite(uid string, contentId string) *Favorite {
	return &Favorite{
		Uid:       uid,
		ContentId: contentId,
		Timestamp: time.Now(),
	}
}
