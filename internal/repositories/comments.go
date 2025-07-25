package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

// Get comments by articleID. Returns slice of comment objects if found, nil otherwise.
func (r *CommentsRepository) GetCommentsByArticleID(articleID string, uid string) ([]models.Comment, error) {
	query := fmt.Sprintf(`
	SELECT
		C.COMMENT_ID,
		C.AUTHOR,
		C.ARTICLE_ID,
		C.CONTENT,
		C.TIME_CREATED,
		C.LAST_EDITED,
		C.FLAGGED,
		C.IS_AVAILABLE,
		C.AUTHOR=$1 AS IS_AUTHOR,
		U.NAME AS AUTHOR_NAME,
		COALESCE(L.LIKE_COUNT, 0) AS NUM_LIKES,
		COALESCE(UL.USER_LIKED, FALSE) AS IS_LIKED
	FROM
		%s C
		INNER JOIN %s U ON C.AUTHOR = U.UID
		LEFT JOIN (
			SELECT
				CONTENT_ID,
				COUNT(*) AS LIKE_COUNT
			FROM
				LIKES
			GROUP BY
				CONTENT_ID
		) L ON L.CONTENT_ID = C.COMMENT_ID
		LEFT JOIN (
			SELECT
				CONTENT_ID,
				TRUE AS USER_LIKED
			FROM
				LIKES
			WHERE
				UID = $1 -- User ID parameter
		) UL ON UL.CONTENT_ID = C.COMMENT_ID
	WHERE
		C.ARTICLE_ID = $2 -- Article ID parameter
		AND C.IS_AVAILABLE = TRUE
	ORDER BY
		C.TIME_CREATED ASC`, COMMENTS_TABLE, USERS_TABLE)

	utils.Logger.Trace().Msgf("Retrieving comments with article_id: %s", articleID)

	rows, _ := r.Db.Query(context.Background(), query, uid, articleID)
	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Comment])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error serializing rows to comment structs")
		return nil, err
	}

	utils.Logger.Trace().Msgf("Found %d comments with article_id %s", len(comments), articleID)

	return comments, nil
}

// Get comment's author
func (r *CommentsRepository) GetAuthor(commentID string) (string, error) {
	query := fmt.Sprintf(`SELECT AUTHOR FROM %s WHERE COMMENT_ID=$1;`, COMMENTS_TABLE)

	utils.Logger.Trace().Msgf("Getting author of comment with id %s", commentID)

	var author string // UID of comment's author
	if err := r.Db.QueryRow(context.Background(), query, commentID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error retrieving author for comment id %s", commentID)
		return "", err
	}

	utils.Logger.Debug().Msgf("Author of comment with id %s is %s", commentID, author)
	return author, nil
}

// Returns true if comment exists in database
func (r *CommentsRepository) IsAvailable(commentID string) bool {
	query := fmt.Sprintf(`SELECT IS_AVAILABLE FROM %s WHERE COMMENT_ID=$1;`, COMMENTS_TABLE)

	var isAvailable bool
	if err := r.Db.QueryRow(context.Background(), query, commentID).Scan(&isAvailable); err != nil {
		utils.Logger.Warn().Msgf("Comment of ID %s not found", commentID)
		return false
	}

	return isAvailable
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
