package services

import (
	"github.com/gofrs/uuid"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type PostService struct {
	PostRepo *repositories.PostsRepository
}

func (s *PostService) CreateNewPost(author string, replyTo *string, threadId string, title string, content string) error {
	post := &models.NewPost{
		Author:   author,
		ThreadId: threadId,
		Title:    title,
		Content:  content,
		ReplyTo:  replyTo,
	}

	return s.PostRepo.CreatePost(post)
}

// Delete post only if author matches the author of the post or if user is admin or staff
func (s *PostService) DeletePost(postId uuid.UUID, claim *utils.JwtClaim) error {
	if claim.Role != "admin" && claim.Role != "staff" {
		author, err := s.PostRepo.GetPostAuthor(postId)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of post")
			return err
		}

		// Check if author of post matches the author in JWT claim
		if author != claim.Username {
			return utils.ErrUnauthorized{}
		}
	}
	return s.PostRepo.DeletePost(postId)
}
