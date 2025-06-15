package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
)

type LikeService struct {
	likesRepository *repositories.LikesRepository
}

var Likes *LikeService

// Create a new like for uid and contentID
func (s *LikeService) CreateLike(uid string, contentID string) error {
	like := models.NewLike(uid, contentID)

	return s.likesRepository.Insert(like)
}

// Check if uid has liked a content
func (s *LikeService) HasLiked(uid string, contentID string) bool {
	return s.likesRepository.GetByUidAndContentId(uid, contentID)
}

// Get number of likes for a content
func (s *LikeService) GetNumLikes(contentID string) int {
	return s.likesRepository.GetNumLikes(contentID)
}

// Remove like for uid and contentID
func (s *LikeService) RemoveLike(uid string, contentID string) error {
	return s.likesRepository.Delete(uid, contentID)
}
