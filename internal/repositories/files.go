package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const FILES_TABLE = "files"

type FilesRepository struct {
	db *pgxpool.Pool
}

var Files *FilesRepository

// Insert file metadata into table. Returns nil on success
func (r *FilesRepository) Insert(file models.DbFile) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (FILE_ID, AUTHOR, FILENAME, GCS_FILENAME, STATUS, FILEGROUP) VALUES $1, $2, $3, $4, $5, $6;`, FILES_TABLE)

	if _, err := r.db.Exec(context.Background(), query, file.FileId, file.AuthorUid, file.Filename, file.GCSFilename, file.Status, file.FileGroup); err != nil {
		utils.Logger.Error().Err(err).Msgf("Failed to insert file %s into postgres", file.FileId)
		return err
	}

	utils.Logger.Info().Interface("file", file).Msg("Successfully inserted file into database")
	return nil
}
