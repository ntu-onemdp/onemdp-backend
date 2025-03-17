package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
)

type LikeService struct {
	likesRepository *repositories.LikesRepository
}

func NewLikeService(likeRepository *repositories.LikesRepository) *LikeService {
	return &LikeService{
		likesRepository: likeRepository,
	}
}

// Create a new like for username and contentID
func (s *LikeService) CreateLike(username string, contentID string) error {
	like := models.NewLike(username, contentID)

	return s.likesRepository.CreateLike(like)
}

// Check if username has liked a content
func (s *LikeService) HasLiked(username string, contentID string) bool {
	return s.likesRepository.GetLikeByUsernameAndContentId(username, contentID)
}

// Get number of likes for a content
func (s *LikeService) GetNumLikes(contentID string) (int, error) {
	return s.likesRepository.GetNumLikes(contentID)
}

// Remove like for username and contentID
func (s *LikeService) RemoveLike(username string, contentID string) error {
	return s.likesRepository.RemoveLike(username, contentID)
}
