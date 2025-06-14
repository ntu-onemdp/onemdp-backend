package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type PostService struct {
	postRepo    *repositories.PostsRepository
	postFactory *models.PostFactory
}

var Posts *PostService

func NewPostService(postRepo *repositories.PostsRepository) *PostService {
	return &PostService{
		postRepo:    postRepo,
		postFactory: models.NewPostFactory(),
	}
}

func (s *PostService) CreateNewPost(author string, replyTo *string, threadId string, title string, content string) error {
	post := s.postFactory.New(author, threadId, title, content, replyTo, false)

	err := s.postRepo.Create(post)
	return err
}

// Retrieve post by post_id
func (s *PostService) GetPost(postID string) (*models.Post, error) {
	return s.postRepo.Get(postID)
}

// Check if post is available (exists and not deleted)
func (s *PostService) PostExists(postID string) bool {
	return s.postRepo.IsAvailable(postID)
}

// Update post. Only the content and the title can be updated.
// Post can only be updated by the author of the post or by admin or staff
func (s *PostService) UpdatePost(updated_post models.Post, uid string) error {
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking staff permission")
		return err
	}

	if !hasStaffPermission {
		author, err := s.postRepo.GetAuthor(updated_post.PostID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of post")
			return err
		}

		// Check if author of post matches the author in JWT claim
		if author != uid {
			return utils.NewErrUnauthorized()
		}
	}

	return s.postRepo.Update(updated_post.PostID, updated_post)
}

// Delete post only if author matches the author of the post or if user is admin or staff
func (s *PostService) DeletePost(postID string, uid string) error {
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking staff permission")
		return err
	}

	if !hasStaffPermission {
		author, err := s.postRepo.GetAuthor(postID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of post")
			return err
		}

		// Check if author of post matches the author in JWT claim
		if author != uid {
			return utils.NewErrUnauthorized()
		}
	}
	return s.postRepo.Delete(postID)
}
