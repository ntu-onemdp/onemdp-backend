package repositories

import "github.com/ntu-onemdp/onemdp-backend/internal/models"

type ContentRepository interface {
	// Insert new content into the database
	Create(content *models.Content) error
	// Get content by content_id. Returns content object if found, nil otherwise.
	GetByID(contentID string) (*models.Content, error)
	// Get content author
	GetAuthor(contentID string) (string, error)
	// Return true if content is available (exists) in database
	IsAvailable(contentID string) bool
	// Update content in the database
	Update(contentID string, content *models.Content) error
	// Update content last updated in the database. Relevant for threads and articles
	UpdateActivity(contentID string) error
	// Soft delete content in the database
	Delete(contentID string) error
	// Restore content in the database
	Restore(contentID string) error
}
