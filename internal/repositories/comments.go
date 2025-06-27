package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Comments table name in db
const COMMENTS_TABLE = "comments"

type CommentsRepository struct {
	_  ContentRepository
	Db *pgxpool.Pool
}

var Comments *CommentsRepository

func (r *CommentsRepository) Create(comment *models.DbComment) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	utils.Logger.Trace().Interface("comment", comment).Msgf("Transaction begin. Inserting comment with id %s", comment.CommentID)

	// Insert comment into comments table
	query := fmt.Sprintf(`
	INSERT INTO %s (comment_id, author, article_id, content)
	VALUES ($1, $2, $3, $4);`, COMMENTS_TABLE)

	if _, err = tx.Exec(ctx, query, comment.CommentID, comment.AuthorUID, comment.ArticleID, comment.Content); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting comment into database")
		return err
	}
	utils.Logger.Trace().Msgf("Comment with id %s successfully inserted into database", comment.CommentID)

	// Update user's karma
	query = fmt.Sprintf(`UPDATE %s SET karma = karma + %d WHERE uid = $1;`, USERS_TABLE, models.COMMENT_ARTICLE_PTS)

	if _, err = tx.Exec(ctx, query, comment.AuthorUID); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating user's karma")
		return err
	}
	utils.Logger.Trace().Msgf("User %s karma successfully updated", comment.AuthorUID)

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Info().Msgf("Comment by %s successfully inserted", comment.AuthorUID)

	return nil
}

// Delete comment from database
func (r *CommentsRepository) Delete(commentID string) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)
	utils.Logger.Trace().Msgf("Transaction begin. Deleting comment with id %s", commentID)

	// Remove comment from comments table and retrieve author uid
	query := fmt.Sprintf(`
	UPDATE %s SET IS_AVAILABLE=false WHERE comment_id=$1 AND IS_AVAILABLE=true RETURNING AUTHOR;
	`, COMMENTS_TABLE)

	var author string // UID of author
	if err = tx.QueryRow(ctx, query, commentID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting comment from database")
		return err
	}

	utils.Logger.Trace().Msgf("Comment with id %s successfully deleted from database", commentID)

	// Remove comment from likes table
	query = fmt.Sprintf(`
	DELETE FROM %s WHERE CONTENT_ID=$1
	`, LIKES_TABLE)

	if _, err = tx.Exec(ctx, query, commentID); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting comments from likes table")
		return err
	}

	utils.Logger.Trace().Msgf("Comments with id %s deleted from likes table", commentID)

	// Update user's karma
	query = fmt.Sprintf(`
	UPDATE %s SET KARMA = GREATEST(KARMA - %d, 0) WHERE UID=$1`, USERS_TABLE, models.COMMENT_ARTICLE_PTS)

	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating user karma")
		return err
	}

	utils.Logger.Trace().Msg("User karma succeessfully updated")

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Info().Msgf("Comment with id %s successfully deleted from database", commentID)
	return nil
}
