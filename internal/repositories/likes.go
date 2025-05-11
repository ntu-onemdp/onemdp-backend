package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Likes table name in db
const LIKES_TABLE = "likes"

type LikesRepository struct {
	Db *pgxpool.Pool
}

var Likes *LikesRepository

// Insert new like into database. Retuns nil on successful insert
func (r *LikesRepository) Insert(like *models.Like) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (username, content_id) 
	VALUES ($1, $2)
	ON CONFLICT (username, content_id) DO NOTHING;`, LIKES_TABLE)

	_, err := r.Db.Exec(context.Background(), query, like.Username, like.ContentId)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return err
	}
	utils.Logger.Info().Msg("Like successfully inserted into database")

	return nil
}

// Get like by username and content_id. Returns true if like exists, false otherwise.
func (r *LikesRepository) GetByUsernameAndContentId(username string, content_id string) bool {
	query := fmt.Sprintf(`SELECT 1 FROM %s WHERE username = $1 AND content_id = $2;`, LIKES_TABLE)

	var num_likes int
	err := r.Db.QueryRow(context.Background(), query, username, content_id).Scan(&num_likes)
	if num_likes == 0 || err != nil {
		utils.Logger.Trace().Str("username", username).Str("content_id", content_id).Msg("Like not found")
		return false
	}

	utils.Logger.Trace().Str("username", username).Str("content_id", content_id).Msg("Like found")
	return true
}

// Get number of likes
func (r *LikesRepository) GetNumLikes(content_id string) int {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE content_id = $1;`, LIKES_TABLE)

	row := r.Db.QueryRow(context.Background(), query, content_id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		utils.Logger.Trace().Err(err).Msg("No likes found")
		return 0
	}

	utils.Logger.Trace().Int("count", count).Msg("Number of likes retrieved")
	return count
}

// Remove like. We perform hard deletes for likes.
func (r *LikesRepository) Delete(username string, content_id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE username = $1 AND content_id = $2;`, LIKES_TABLE)

	_, err := r.Db.Exec(context.Background(), query, username, content_id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting like")
		return err
	}

	utils.Logger.Info().Msg("Like successfully deleted")
	return nil
}
