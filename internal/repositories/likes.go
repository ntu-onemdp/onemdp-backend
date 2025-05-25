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
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Insert like into database
	query := fmt.Sprintf(`	
	INSERT INTO %s (username, content_id) 
	VALUES ($1, $2)
	ON CONFLICT (username, content_id) DO NOTHING;
	`, LIKES_TABLE)

	if _, err = tx.Exec(context.Background(), query, like.Username, like.ContentId); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting like into database")
		return err
	}

	// Get username of author
	contentType := string(like.ContentId[0]) // Content type is the first character of content_id
	var author string                        // Username of author

	switch contentType {
	case "p": // Post
		query = `SELECT author FROM posts WHERE post_id = $1;`
	case "t": // Thread
		query = `SELECT author FROM threads WHERE thread_id = $1;`
	default:
		utils.Logger.Error().Msg("Unknown content type")
		return fmt.Errorf("unknown content type: %s", contentType)
	}

	if err = tx.QueryRow(ctx, query, like.ContentId).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting author of content")
		return err
	}

	// Update author's karma
	query = fmt.Sprintf(`UPDATE %s SET karma = karma + %d WHERE username = $1;`, USERS_TABLE, models.LIKE_PTS)
	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma")
		return err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
	}

	utils.Logger.Debug().Msg("Like successfully inserted into database")
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
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Delete like from database
	query := fmt.Sprintf(`DELETE FROM %s WHERE username = $1 AND content_id = $2;`, LIKES_TABLE)
	if _, err = tx.Exec(ctx, query, username, content_id); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting like from database")
		return err
	}

	// Get username of author
	contentType := string(content_id[0]) // Content type is the first character of content_id
	var author string                    // Username of author
	switch contentType {
	case "p": // Post
		query = `SELECT author FROM posts WHERE post_id = $1;`
	case "t": // Thread
		query = `SELECT author FROM threads WHERE thread_id = $1;`
	default:
		utils.Logger.Error().Msg("Unknown content type")
		return fmt.Errorf("unknown content type: %s", contentType)
	}

	if err = tx.QueryRow(ctx, query, content_id).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting author of content")
		return err
	}

	// Update author's karma
	query = fmt.Sprintf(`UPDATE %s SET karma = GREATEST(karma - %d, 0) WHERE username = $1;`, USERS_TABLE, models.LIKE_PTS)
	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma")
		return err
	}
	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Info().Msg("Like successfully deleted")
	return nil
}
