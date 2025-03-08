package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
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

// Utility function to get preview from content
func getPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 100

	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}
