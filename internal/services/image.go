package services

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"

	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ImageService struct {
	imageRepo *repositories.ImagesRepository
}

var Images *ImageService

// Retrieve image from database
func (s *ImageService) Get(id string) ([]byte, error) {
	return s.imageRepo.Get(id)
}

// Process image into bytes array and pass it to the repository. Returns UUID of image on success
func (s *ImageService) Insert(image *multipart.FileHeader) (string, error) {
	// Reject is image size is too large
	if image.Size > 2*1024*1024 { // 2 MB
		utils.Logger.Error().Msg("Image size exceeds 2 MB")
		return "", errors.New("image size exceeds 2 MB")
	}

	// Reject if image type is not supported
	if !isValidType(image) {
		utils.Logger.Error().Msg("Unsupported image type")
		return "", errors.New("unsupported image type")
	}

	bytes := make([]byte, image.Size)

	// Open image file
	file, err := image.Open()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to open image file")
		return "", err
	}
	defer file.Close()

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

// Check if image is of supported type
func isValidType(image *multipart.FileHeader) bool {
	imgType := image.Header.Get("Content-Type")

	switch imgType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	}
	return false
}
