package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadService struct {
	ThreadRepo *repositories.ThreadRepository
	PostRepo   *repositories.PostsRepository
}

// Create new thread and insert into the repository
func (s *ThreadService) CreateNewThread(author string, title string, content string) error {
	thread := &models.NewThread{
		Author:  author,
		Title:   title,
		Preview: getPreview(content),
	}
	threadId, err := s.ThreadRepo.CreateThread(thread)
	if err != nil {
		return err
	}

	post := &models.NewPost{
		Author:   author,
		ThreadId: threadId,
		Title:    title,
		Content:  content,
		ReplyTo:  nil,
	}

	return s.PostRepo.CreatePost(post)
}

// Retrieve thread and all associated posts
func (s *ThreadService) GetThread(threadId string) (*models.Thread, []models.Post, error) {
	thread, err := s.ThreadRepo.GetThreadById(threadId)
	if err != nil {
		return nil, nil, err
	}

	posts, err := s.PostRepo.GetPostByThreadId(threadId)
	if err != nil {
		return nil, nil, err
	}

	return thread, posts, nil
}

// Delete thread and all associated posts
func (s *ThreadService) DeleteThread(threadId string, claim *utils.JwtClaim) error {
	// Check if role is admin or staff
	if !HasStaffPermission(claim) {
		author, err := s.ThreadRepo.GetThreadAuthor(threadId)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != claim.Username {
			return utils.NewErrUnauthorized()
		}
	}

	err := s.ThreadRepo.DeleteThread(threadId)
	if err != nil {
		utils.Logger.Trace().Msg("Error deleting thread")
		return err
	}

	return s.PostRepo.DeletePostsByThread(threadId)
}

// Utility function to get preview from content
func getPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 100

	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}
