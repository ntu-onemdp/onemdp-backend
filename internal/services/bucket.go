// Manage connections to Google Cloud Storage
package services

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/ntu-onemdp/onemdp-backend/internal/db"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type GCSFileService struct {
	bucket *storage.BucketHandle
	dir    string // Directory to save the files in
}

func NewGCSFileService() *GCSFileService {
	dir, found := os.LookupEnv("GCS_DIR")
	if !found {
		utils.Logger.Warn().Msg("GCS_DIR not set in environment variables")
	}

	return &GCSFileService{
		bucket: db.Bucket,
		dir:    dir,
	}
}

var GCSFileServiceInstance *GCSFileService

func init() {
	GCSFileServiceInstance = NewGCSFileService()
}

// Upload file to pdfstore in GCS
func (s *GCSFileService) Upload(file *multipart.FileHeader, filename string) error {
	handler := s.bucket.Object(s.dir + filename)
	utils.Logger.Trace().Msg("handle to bucket object created")

	writer := handler.NewWriter(context.Background())

	f, err := file.Open()
	if err != nil {
		utils.Logger.Error().Err(err).Str("filename", file.Filename).Msg("Error encountered when opening file")
		return err
	}
	defer f.Close()

	// Stream file contents into the bucket's object
	if _, err := io.Copy(writer, f); err != nil {
		utils.Logger.Error().Err(err).Str("filename", file.Filename).Msg("Error encountered during file upload to GCS bucket")
		writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		utils.Logger.Error().Err(err).Msg("Error closing writer object")
		return err
	}

	utils.Logger.Info().Str("gcs filename", filename).Msg("File successfully uploaded to GCS")
	return nil
}

// Retrieve file from GCS
