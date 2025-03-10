package repositories

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type PostsRepository struct {
	Db *pgxpool.Pool
}

// Posts table name in db
const POSTS_TABLE = "posts"

// Insert new post into the database. Returns nil on successful insert
func (r *PostsRepository) CreatePost(post *models.NewPost) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (author, thread_id, title, content, reply_to) 
	VALUES ($1, $2, $3, $4, $5);`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Inserting post: %v", post))

	_, err := r.Db.Exec(context.Background(), query, post.Author, post.ThreadId, post.Title, post.Content, post.ReplyTo)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully inserted into database", post.Title))
	return nil
}

// Get post author
func (r *PostsRepository) GetPostAuthor(postId uuid.UUID) (string, error) {
	query := fmt.Sprintf(`SELECT author FROM %s WHERE post_id = $1;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Getting author of post with id: %v", postId))

	var author string
	err := r.Db.QueryRow(context.Background(), query, postId).Scan(&author)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return "", err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Author of post with id %v is %v", postId, author))
	return author, nil
}

// Delete post from database using post_id. Returns nil if successful.
// Soft delete is performed.
func (r *PostsRepository) DeletePost(postId uuid.UUID) error {
	query := fmt.Sprintf(`
	WITH deleted_rows AS (
		UPDATE %s SET is_available = false WHERE post_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM deleted_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Deleting post with id: %v", postId))

	var num_deleted int
	err := r.Db.QueryRow(context.Background(), query, postId).Scan(&num_deleted)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	if num_deleted == 0 {
		utils.Logger.Warn().Msg("No rows deleted")
	}

	utils.Logger.Info().Int("Rows deleted", num_deleted).Msg(fmt.Sprintf("Post with id %v successfully deleted from database", postId))
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
func (r *PostsRepository) RestorePost(postId uuid.UUID) error {
	query := fmt.Sprintf(`
	WITH restored_rows AS (
		UPDATE %s SET is_available = true WHERE post_id = $1 RETURNING thread_id
	)
	SELECT COUNT(*) FROM restored_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Restoring post with id: %v", postId))

	var num_restored int
	err := r.Db.QueryRow(context.Background(), query, postId).Scan(&num_restored)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring post")
		return err
	}

	utils.Logger.Info().Int("Rows restored", num_restored).Msg(fmt.Sprintf("Post with id %v successfully restored", postId))
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
