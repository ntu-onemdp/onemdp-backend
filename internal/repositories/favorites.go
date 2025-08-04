package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const FAVORITES_TABLE = "favorites"

type FavoritesRepository struct {
	Db *pgxpool.Pool
}

var Favorites *FavoritesRepository

// Insert new favorite into the database, return nil on successful insert.
func (r *FavoritesRepository) Insert(favorite *models.Favorite) error {
	ctx := context.Background()

	query := fmt.Sprintf(`INSERT INTO %s (UID, CONTENT_ID) VALUES ($1, $2) ON CONFLICT (UID, CONTENT_ID) DO NOTHING;`, FAVORITES_TABLE)

	if _, err := r.Db.Exec(ctx, query, favorite.Uid, favorite.ContentId); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting favorite into database for %s", favorite.Uid)
		return err
	}

	utils.Logger.Debug().Str("uid", favorite.Uid).Str("content id", favorite.ContentId).Msgf("Successfully inserted into favorites table for %s", favorite.Uid)

	return nil
}

// Get favorite by uid and content_id, return true if record exists in database and false otherwise.
func (r *FavoritesRepository) Exists(uid string, contentID string) bool {
	query := fmt.Sprintf(`SELECT 1 FROM %s WHERE uid = $1 AND content_id = $2;`, FAVORITES_TABLE)

	var numRecords int
	err := r.Db.QueryRow(context.Background(), query, uid, contentID).Scan(&numRecords)
	if numRecords == 0 || err != nil {
		utils.Logger.Trace().Str("uid", uid).Str("content_id", contentID).Msg("Record not found")
		return false
	}

	utils.Logger.Trace().Str("uid", uid).Str("content_id", contentID).Msg("Record found")
	return true
}

// Remove favorite from database (hard delete is performed)
// Returns nil on success
func (r *FavoritesRepository) Delete(uid string, contentID string) error {
	ctx := context.Background()

	query := fmt.Sprintf(`DELETE FROM %s WHERE UID=$1 AND CONTENT_ID=$2;`, FAVORITES_TABLE)

	if _, err := r.Db.Exec(ctx, query, uid, contentID); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting favorite into database for %s", uid)
		return err
	}

	utils.Logger.Debug().Str("uid", uid).Str("content id", contentID).Msgf("Successfully removed from favorites table for %s", uid)

	return nil
}
