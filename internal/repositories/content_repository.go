package repositories

import "github.com/ntu-onemdp/onemdp-backend/internal/models"

type ContentRepository interface {
	Create(content *models.Content) error
	GetByID(contentID string) (*models.Content, error)
	GetAuthor(contentID string) (string, error)
	IsAvailable(contentID string) bool
	Update(contentID string, content *models.Content) error
	UpdateActivity(contentID string) error
	Delete(contentID string) error
	Restore(contentID string) error
}
