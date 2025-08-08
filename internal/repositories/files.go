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
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Replace old files.
	query := fmt.Sprintf(`UPDATE %s SET STATUS='available', TIME_DELETED=NOW(), DELETED_BY=$1 WHERE FILENAME=$2 AND FILE_GROUP=$3;`, FILES_TABLE)
	if _, err := tx.Exec(ctx, query, file.AuthorUid, file.Filename, file.FileGroup); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating metadata for old files")
		return err
	}

	// Insert new file
	query = fmt.Sprintf(`
	INSERT INTO %s (FILE_ID, AUTHOR, FILENAME, GCS_FILENAME, STATUS, FILE_GROUP) VALUES ($1, $2, $3, $4, $5, $6);`, FILES_TABLE)

	if _, err := tx.Exec(ctx, query, file.FileId, file.AuthorUid, file.Filename, file.GCSFilename, file.Status, file.FileGroup); err != nil {
		utils.Logger.Error().Err(err).Msgf("Failed to insert file %s into postgres", file.FileId)
		return err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Info().Interface("file", file).Msg("Successfully inserted file into database")
	return nil
}
