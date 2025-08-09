package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	query := fmt.Sprintf(`UPDATE %s SET STATUS='deleted', TIME_DELETED=NOW(), DELETED_BY=$1 WHERE FILENAME=$2 AND FILE_GROUP=$3;`, FILES_TABLE)
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

// Get GCS filename and original filename from id
// Note that other fields are not retrieved. Accessing them will net default values.
func (r *FilesRepository) GetFilename(id string) (*models.DbFile, error) {
	query := fmt.Sprintf(`SELECT GCS_FILENAME, FILENAME FROM %s WHERE FILE_ID=$1 AND STATUS='available';`, FILES_TABLE)

	var filename *string
	var GCSFilename *string
	if err := r.db.QueryRow(context.Background(), query, id).Scan(&GCSFilename, &filename); err != nil {
		if err == pgx.ErrNoRows {
			utils.Logger.Error().Err(err).Str("file id", id).Msgf("File with %s does not exist", id)
			return nil, err
		}

		utils.Logger.Error().Err(err).Str("file id", id).Msgf("Error fetching GCS filename for id %s", id)
		return nil, err
	}

	utils.Logger.Debug().Str("file id", id).Str("GCS filename", *filename).Msgf("GCS filename for file %s found", id)
	return &models.DbFile{Filename: *filename, GCSFilename: *GCSFilename}, nil
}

// Get list of files available
func (r *FilesRepository) GetFileList() ([]models.FileMetadata, error) {
	query := fmt.Sprintf(`
	SELECT 
		F.FILE_ID, F.AUTHOR, F.FILENAME, F.TIME_CREATED, F.FILE_GROUP, U.NAME AS AUTHOR_NAME
	FROM 
		%s F-- Files table
		INNER JOIN %s U ON F.AUTHOR=U.UID -- Users table
	WHERE 
		F.STATUS='available';`, FILES_TABLE, USERS_TABLE)

	rows, _ := r.db.Query(context.Background(), query)
	files, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.FileMetadata])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows into struct")
		return nil, err
	}

	utils.Logger.Debug().Msgf("%d files retrieved from postgres", len(files))
	return files, nil
}

// Revert change if upload to GCS bucket is unsuccessful
// Params: fileID of file to remove from database
func (r *FilesRepository) Revert(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE FILE_ID=$1;`, FILES_TABLE)

	if _, err := r.db.Exec(context.Background(), query, id); err != nil {
		utils.Logger.Error().Err(err).Msg("Error reverting change")
		return err
	}

	utils.Logger.Info().Str("File ID", id).Msgf("File ID %s removed from databse", id)
	return nil
}

// Perform soft delete of file in database
func (r *FilesRepository) Delete(id string, uid string) error {
	query := fmt.Sprintf(`
	UPDATE %s 
	SET STATUS='deleted', DELETED_BY=$1, TIME_DELETED=NOW()
	WHERE FILE_ID=$2;`, FILES_TABLE)

	if _, err := r.db.Exec(context.Background(), query, uid, id); err != nil {
		utils.Logger.Error().Err(err).Str("uid", uid).Str("file id", id).Msg("Error deleting file from database")
		return err
	}

	utils.Logger.Info().Str("uid", uid).Str("file id", id).Msg("File metadata successfully removed from database")
	return nil
}
