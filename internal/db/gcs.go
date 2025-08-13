package db

import (
	"cloud.google.com/go/storage"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

var Bucket *storage.BucketHandle

func init() {
	// env := os.Getenv("ENV")
	// var path string
	// if env == "DEV" {
	// 	path = "secrets/service-account-key"
	// } else {
	// 	path = "mnt/secrets/service-account-key"
	// }

	// client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(path))
	// if err != nil {
	// 	utils.Logger.Panic().Err(err).Msg("Error connecting to Google Cloud Storage.")
	// }
	// defer client.Close()

	// name, found := os.LookupEnv("GCS_BUCKET_NAME")
	// if !found {
	// 	utils.Logger.Warn().Msg("GCS_BUCKET_NAME not set in environment variables.")
	// }
	// Bucket = client.Bucket(name)

	utils.Logger.Info().Msg("Connection to GCS bucket established.")
}
