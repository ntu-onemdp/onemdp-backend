package services

import (
	"errors"
	"fmt"
	"mime/multipart"

	c "github.com/ntu-onemdp/onemdp-backend/config"
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
	if image.Size > c.MAX_IMAGE_SIZE {
		err := fmt.Sprintf("image size exceeds %d MB", c.MAX_IMAGE_SIZE/(1024*1024))
		utils.Logger.Error().Msg(err)
		return "", errors.New(err)
	}

	// Reject if image type is not supported
	if !isValidType(image) {
		utils.Logger.Error().Msg("Unsupported image type")
		return "", errors.New("unsupported image type")
	}

	// Sanitize image
	img, err := utils.SanitizeImage(image)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to sanitize image")
		return "", err
	}

	// Insert into db
	id, err := s.imageRepo.Insert(img)
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
