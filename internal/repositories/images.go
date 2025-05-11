package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const IMAGES_TABLE = "images"

type ImagesRepository struct {
	Db *pgxpool.Pool
}

var Images *ImagesRepository

// Insert new image into the database.Returns UUID of image on success
