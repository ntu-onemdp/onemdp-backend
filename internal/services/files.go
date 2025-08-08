package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
)

type FileService struct {
	fileRepo *repositories.FilesRepository

	fileFactory *models.FileFactory
}

func NewFileService(fileRepo *repositories.FilesRepository) *FileService {
	return &FileService{
		fileRepo:    fileRepo,
		fileFactory: models.NewFileFactory(),
	}
}

var Files *FileService

func (s *FileService) Create(author string, filename string, filegroup *string) error {
	dbFile := s.fileFactory.New(author, filename, filegroup)

	return s.fileRepo.Insert(*dbFile)
}
