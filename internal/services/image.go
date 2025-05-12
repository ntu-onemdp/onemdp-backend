package services

import (
	"bufio"
	"io"
	"mime/multipart"

	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ImageService struct {
	imageRepo *repositories.ImagesRepository
}

var Images *ImageService

// Process image into bytes array and pass it to the repository. Returns UUID of image on success
func (s *ImageService) Insert(image *multipart.FileHeader) (string, error) {
	bytes := make([]byte, image.Size)

	// Open image file
	file, err := image.Open()
	defer file.Close()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to open image file")
		return "", err
	}

	// Read image file into bytes
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil && err != io.EOF {
		utils.Logger.Error().Err(err).Msg("Failed to read image file")
		return "", err
	}

	// Insert into db
	id, err := s.imageRepo.Insert(bytes)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to insert image into db")
		return "", err
	}

	return id, nil
}
