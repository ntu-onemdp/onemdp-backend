package services

import (
	"github.com/gofrs/uuid"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
)

type PostService struct {
	PostRepo *repositories.PostsRepository
}

func (s *PostService) CreateNewPost(author string, replyTo *string, threadId uuid.UUID, title string, content string) error {
	post := &models.NewPost{
		Author:   author,
		ThreadId: threadId,
		Title:    title,
		Content:  content,
		ReplyTo:  replyTo,
	}

	return s.PostRepo.CreatePost(post)
}
