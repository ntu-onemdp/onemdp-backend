package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const IMAGES_TABLE = "images"

type ImagesRepository struct {
	Db *pgxpool.Pool
}

var Images *ImagesRepository

// Insert new image into the database.Returns UUID of image on success
func (r *ImagesRepository) Insert(image []byte) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s (image) VALUES ($1) RETURNING id;`, IMAGES_TABLE)

	var id string
	err := r.Db.QueryRow(context.Background(), query, image).Scan(&id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to insert image into db")
		return "", err
	}

	utils.Logger.Debug().Msgf("Inserted image with ID: %s", id)
	return id, nil
}
