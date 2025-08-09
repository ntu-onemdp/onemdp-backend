package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
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

// Retrieve GCS filename and original filename from DB
// Note that other fields are not retrieved. Accessing them will net default values.
func (s *FileService) GetFilename(id string) (*models.DbFile, error) {
	return s.fileRepo.GetFilename(id)
}

// Get list of files available and group them
func (s *FileService) GetFileList() (map[string][]models.FileMetadata, error) {
	filelist, err := s.fileRepo.GetFileList()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error retrieving list from database")
		return nil, err
	}

	// Group it by file group
	list := make(map[string][]models.FileMetadata, len(filelist))
	for _, f := range filelist {
		list[*f.FileGroup] = append(list[*f.FileGroup], f)
	}

	return list, nil
}

// Revert change if upload to GCS bucket is unsuccessful
func (s *FileService) Revert(id string) error {
	return s.fileRepo.Revert(id)
}

// Delete file from postgres
func (s *FileService) Remove(id string, uid string) error {
	return s.fileRepo.Delete(id, uid)
}
