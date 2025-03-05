package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
)

type ThreadService struct {
	ThreadRepo *repositories.ThreadRepository
}

// Create new thread and insert into the repository
func (s *ThreadService) CreateNewThread(author string, title string, content string) error {
	thread := &models.NewThread{
		Author:  author,
		Title:   title,
		Preview: content,
	}
	return s.ThreadRepo.CreateThread(thread)
}
