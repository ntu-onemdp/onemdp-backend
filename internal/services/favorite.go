package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type FavoriteService struct {
	repo *repositories.FavoritesRepository
}

var Favorites *FavoriteService

// Favorite a new content
func (s *FavoriteService) CreateFavorite(uid string, contentID string) error {
	// Raise warning if favorite is of unsupported content type
	if string(contentID[0]) != "a" && string(contentID[0]) != "t" {
		utils.Logger.Warn().Str("content id", contentID).Msg("Unsupported content id, you should only favorite threads and articles.")
	}

	favorite := models.NewFavorite(uid, contentID)
	return s.repo.Insert(favorite)
}

func (s *FavoriteService) Exists(uid string, contentID string) bool {
	return s.repo.Exists(uid, contentID)
}

// Remove favorite
func (s *FavoriteService) RemoveFavorite(uid string, contentID string) error {
	return s.repo.Delete(uid, contentID)
}
