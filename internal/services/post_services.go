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

// Update post. Only the content and the title can be updated.
// Post can only be updated by the author of the post or by admin or staff
func (s *PostService) UpdatePost(updated_post models.Post, claim *utils.JwtClaim) error {
	if !HasStaffPermission(claim) {
		author, err := s.PostRepo.GetPostAuthor(updated_post.PostId)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of post")
			return err
		}

		// Check if author of post matches the author in JWT claim
		if author != claim.Username {
			return utils.ErrUnauthorized{}
		}
	}

	restructured_post := &models.NewPost{
		Author:   updated_post.Author,
		ThreadId: updated_post.ThreadId,
		Title:    updated_post.Title,
		Content:  updated_post.Content,
		ReplyTo:  updated_post.ReplyTo,
	}

	return s.PostRepo.UpdatePostContent(updated_post.PostId, *restructured_post)
}

// Delete post only if author matches the author of the post or if user is admin or staff
func (s *PostService) DeletePost(postId uuid.UUID, claim *utils.JwtClaim) error {
	if !HasStaffPermission(claim) {
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
