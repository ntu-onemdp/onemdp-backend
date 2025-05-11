package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Posts table name in db
const POSTS_TABLE = "posts"

type PostsRepository struct {
	_  ContentRepository
	Db *pgxpool.Pool
}

var Posts *PostsRepository

// Insert new post into the database. Returns post ID and nil on successful insert
func (r *PostsRepository) Create(post *models.Post) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (post_id, author, thread_id, title, content, reply_to, is_header) 
	VALUES ($1, $2, $3, $4, $5, $6, $7);`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Inserting post: %v", post))

	_, err := r.Db.Exec(context.Background(), query, post.PostID, post.Author, post.ThreadId, post.Title, post.PostContent, post.ReplyTo, post.IsHeader)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully inserted into database", post.Title))
	return nil
}

// Get posts by thread_id. Returns slice of post objects if found, nil otherwise.
func (r *PostsRepository) GetPostByThreadId(threadId string) ([]models.Post, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE thread_id = $1 AND is_available = true ORDER BY time_created ASC;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting posts with thread_id: %s", threadId))

	rows, _ := r.Db.Query(context.Background(), query, threadId)
	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Post])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error serializing rows to post structs")
		return nil, err
	}

	utils.Logger.Trace().Interface("Posts", posts).Msg(fmt.Sprintf("Posts with thread_id %s found", threadId))
	return posts, nil
}

// Get number of replies by in a thread. Returns number of replies if found, 0 otherwise.
// Note that the header post is not counted.
func (r *PostsRepository) GetNumReplies(threadId string) int {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE thread_id = $1 AND is_header = false AND is_available = true;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting number of replies with thread_id: %s", threadId))

	var numReplies int
	err := r.Db.QueryRow(context.Background(), query, threadId).Scan(&numReplies)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return 0
	}

	utils.Logger.Info().Int("Number of replies", numReplies).Msg(fmt.Sprintf("Number of replies with thread_id %s found", threadId))
	return numReplies
}

// Get post author
func (r *PostsRepository) GetAuthor(postID string) (string, error) {
	query := fmt.Sprintf(`SELECT author FROM %s WHERE post_id = $1;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting author of post with id: %v", postID))

	var author string
	err := r.Db.QueryRow(context.Background(), query, postID).Scan(&author)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return "", err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Author of post with id %v is %v", postID, author))
	return author, nil
}

// Check if post exists
func (r *PostsRepository) IsAvailable(postID string) bool {
	query := fmt.Sprintf(`SELECT is_available FROM %s WHERE post_id = $1;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Checking if post with id: %v exists", postID))

	var isAvailable bool
	err := r.Db.QueryRow(context.Background(), query, postID).Scan(&isAvailable)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return false
	}

	utils.Logger.Trace().Bool("Is available", isAvailable).Msg(fmt.Sprintf("Post with id %v exists", postID))
	return isAvailable
}

// Edit content of post
func (r *PostsRepository) Update(postID string, updated_post models.Post) error {
	query := fmt.Sprintf(`
	UPDATE %s SET title = $1, content = $2, reply_to = $3, last_edited = NOW() WHERE post_id = $4;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Updating content of post with id: %v", postID))

	_, err := r.Db.Exec(context.Background(), query, updated_post.Title, updated_post.Content, updated_post.ReplyTo, postID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Content of post with id %v successfully updated", postID))
	return nil
}

// Delete post from database using post_id. Returns nil if successful.
// Soft delete is performed.
func (r *PostsRepository) Delete(postID string) error {
	query := fmt.Sprintf(`
	WITH deleted_rows AS (
		UPDATE %s SET is_available = false WHERE post_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM deleted_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Deleting post with id: %v", postID))

	var num_deleted int
	err := r.Db.QueryRow(context.Background(), query, postID).Scan(&num_deleted)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	if num_deleted == 0 {
		utils.Logger.Warn().Msg("No rows deleted")
	}

	utils.Logger.Info().Int("Rows deleted", num_deleted).Msg(fmt.Sprintf("Post with id %v successfully deleted from database", postID))
	return nil
}

// Delete all posts from database matching thread_id. Returns nil if successful.
func (r *PostsRepository) DeletePostsByThread(threadId string) error {
	query := fmt.Sprintf(`
	WITH deleted_rows AS (
		UPDATE %s SET is_available = false WHERE thread_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM deleted_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Deleting posts with thread_id: %s", threadId))

	var num_deleted int
	err := r.Db.QueryRow(context.Background(), query, threadId).Scan(&num_deleted)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	if num_deleted == 0 {
		utils.Logger.Warn().Msg("No rows deleted")
	}

	utils.Logger.Info().Int("Rows deleted", num_deleted).Msg(fmt.Sprintf("Posts with thread_id %s successfully deleted from database", threadId))
	return nil
}

// Restore deleted post by post_id. Returns nil if successful.
func (r *PostsRepository) Restore(postID string) error {
	query := fmt.Sprintf(`
	WITH restored_rows AS (
		UPDATE %s SET is_available = true WHERE post_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM restored_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Restoring post with id: %v", postID))

	var num_restored int
	err := r.Db.QueryRow(context.Background(), query, postID).Scan(&num_restored)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring post")
		return err
	}

	utils.Logger.Info().Int("Rows restored", num_restored).Msg(fmt.Sprintf("Post with id %v successfully restored", postID))
	return nil
}

// Restore all deleted posts by thread_id. Returns nil if successful.
func (r *PostsRepository) RestorePostsByThread(threadId string) error {
	query := fmt.Sprintf(`
	WITH restored_rows AS (
		UPDATE %s SET is_available = true WHERE thread_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM restored_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Restoring posts with thread_id: %s", threadId))

	var num_restored int
	err := r.Db.QueryRow(context.Background(), query, threadId).Scan(&num_restored)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring posts")
		return err
	}

	utils.Logger.Info().Int("Rows restored", num_restored).Msg(fmt.Sprintf("Posts with thread_id %s successfully restored", threadId))
	return nil
}
