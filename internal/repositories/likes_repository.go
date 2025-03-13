package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type LikesRepository struct {
	Db *pgxpool.Pool
}

// Likes table name in db
const LIKES_TABLE = "likes"

// Insert new like into database. Retuns nil on successful insert
func (r *LikesRepository) CreateLike(like *models.Like) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (username, content_id, content_type) 
	VALUES ($1, $2, $3);`, LIKES_TABLE)

	_, err := r.Db.Exec(context.Background(), query, like.Username, like.ContentId, like.ContentType)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return err
	}
	utils.Logger.Info().Msg("Like successfully inserted into database")

	return nil
}

// Get like by username and content_id. Returns true if like exists, false otherwise.
func (r *LikesRepository) GetLikeByUsernameAndContentId(username string, content_id string) (bool, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username = $1 AND content_id = $2;`, LIKES_TABLE)

	row := r.Db.QueryRow(context.Background(), query, username, content_id)
	var like models.Like
	err := row.Scan(&like.Username, &like.ContentId, &like.ContentType)
	if err != nil {
		utils.Logger.Trace().Str("username", username).Str("content_id", content_id).Msg("Like not found")
		return false, nil
	}

	utils.Logger.Trace().Str("username", username).Str("content_id", content_id).Msg("Like found")
	return true, nil
}

// Remove like. We perform hard deletes for likes.
func (r *LikesRepository) RemoveLike(username string, content_id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE username = $1 AND content_id = $2;`, LIKES_TABLE)

	_, err := r.Db.Exec(context.Background(), query, username, content_id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting like")
		return err
	}

	utils.Logger.Info().Msg("Like successfully deleted")
	return nil
}
