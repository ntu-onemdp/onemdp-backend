package db

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
	"google.golang.org/api/option"
)

var Bucket *storage.BucketHandle

func init() {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile("secrets/service-account-key.json"))
	if err != nil {
		utils.Logger.Panic().Err(err).Msg("Error connecting to Google Cloud Storage.")
	}
	defer client.Close()

	name, found := os.LookupEnv("GCS_BUCKET_NAME")
	if !found {
		utils.Logger.Warn().Msg("GCS_BUCKET_NAME not set in environment variables.")
	}
	Bucket = client.Bucket(name)

	utils.Logger.Info().Msg("Connection to GCS bucket established.")
}
