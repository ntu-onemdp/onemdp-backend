package utils

import (
	"bytes"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/microcosm-cc/bluemonday"
)

// Sanitize content to prevent XSS attacks
func SanitizeContent(content string) string {
	policy := bluemonday.UGCPolicy()

	// Allow styles on images (to allow for image resizing)
	policy.AllowStyles("width", "height", "draggable").OnElements("img")

	return policy.Sanitize(content)
}

// Sanitize image
func SanitizeImage(f *multipart.FileHeader) ([]byte, error) {
	file, err := f.Open()
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to open image file")
		return nil, err
	}

	// Decode image
	decoded, err := imaging.Decode(file)
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to decode image")
		return nil, err
	}

	// Encode image to bytes
	format, err := imaging.FormatFromFilename(f.Filename)
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to get image format")
		return nil, err
	}
	var buf bytes.Buffer
	err = imaging.Encode(&buf, decoded, format, imaging.JPEGQuality(80))
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to encode image")
		return nil, err
	}

	return buf.Bytes(), nil
}
