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

func (s *FileService) Create(author string, filename string, filegroup *string) (*models.DbFile, error) {
	dbFile := s.fileFactory.New(author, filename, filegroup)

	return dbFile, s.fileRepo.Insert(*dbFile)
}

// Retrieve GCS filename from DB
func (s *FileService) GetGCSFilename(id string) (string, error) {
	return s.fileRepo.GetGCSFilename(id)
}

// Revert change if upload to GCS bucket is unsuccessful
func (s *FileService) Revert(id string) error {
	return s.fileRepo.Revert(id)
}

// Delete file from postgres
func (s *FileService) Remove(id string, uid string) error {
	return s.fileRepo.Delete(id, uid)
}
