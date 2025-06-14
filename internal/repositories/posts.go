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
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)
	utils.Logger.Trace().Interface("post", post).Msg(fmt.Sprintf("Transaction begin. Inserting post with id: %v", post.PostID))

	// Insert post into posts table
	query := fmt.Sprintf(`
	INSERT INTO %s (post_id, author, thread_id, title, content, reply_to, is_header) 
	VALUES ($1, $2, $3, $4, $5, $6, $7);`, POSTS_TABLE)

	if _, err = tx.Exec(ctx, query, post.PostID, post.Author, post.ThreadId, post.Title, post.PostContent, post.ReplyTo, post.IsHeader); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting post into database")
		return err
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("Post with id %v successfully inserted into database", post.PostID))

	// Update user karma
	// Do not update karma if post is a header post as karma would have been updated when the thread was created
	if !post.IsHeader {
		query = fmt.Sprintf(`
		UPDATE %s SET karma = karma + %d WHERE uid = $1;`, USERS_TABLE, models.CREATE_POST_PTS)

		if _, err = tx.Exec(ctx, query, post.Author); err != nil {
			utils.Logger.Error().Err(err).Msg("Error updating user karma")
			return err
		}
		utils.Logger.Trace().Msg(fmt.Sprintf("User %s karma successfully updated", post.Author))
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("Transaction committed. Post with id %s successfully inserted into database", post.PostID))

	utils.Logger.Debug().Msg(fmt.Sprintf("%s successfully inserted into database", post.Title))
	return nil
}

// Get post by post_id. Returns post object if found, nil otherwise.
func (r *PostsRepository) Get(postID string) (*models.Post, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE post_id = $1 AND is_available = true;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting post with id: %s", postID))

	row, _ := r.Db.Query(context.Background(), query, postID)
	post, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Post])

	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("Post with id %v found", postID))
	return &post, nil
}

// Get posts by thread_id. Returns slice of post objects if found, nil otherwise.
func (r *PostsRepository) GetPostByThreadId(threadID string) ([]models.Post, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE thread_id = $1 AND is_available = true ORDER BY time_created ASC;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting posts with thread_id: %s", threadID))

	rows, _ := r.Db.Query(context.Background(), query, threadID)
	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Post])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error serializing rows to post structs")
		return nil, err
	}

	utils.Logger.Debug().Interface("Posts", posts).Msg(fmt.Sprintf("Posts with thread_id %s found", threadID))
	return posts, nil
}

// Get number of replies by in a thread. Returns number of replies if found, 0 otherwise.
// Note that the header post is not counted.
func (r *PostsRepository) GetNumReplies(threadID string) int {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE thread_id = $1 AND is_header = false AND is_available = true;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting number of replies with thread_id: %s", threadID))

	var numReplies int
	err := r.Db.QueryRow(context.Background(), query, threadID).Scan(&numReplies)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return 0
	}

	utils.Logger.Debug().Int("Number of replies", numReplies).Msg(fmt.Sprintf("Number of replies with thread_id %s found", threadID))
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

	utils.Logger.Debug().Msg(fmt.Sprintf("Author of post with id %v is %v", postID, author))
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

	utils.Logger.Debug().Bool("Is available", isAvailable).Msg(fmt.Sprintf("Post with id %v exists", postID))
	return isAvailable
}

// Edit content of post
// TODO: update last edited for parent thread
func (r *PostsRepository) Update(postID string, updated_post models.Post) error {
	query := fmt.Sprintf(`
	UPDATE %s SET title = $1, content = $2, reply_to = $3, last_edited = NOW() WHERE post_id = $4 AND is_available = true;`, POSTS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Updating content of post with id: %v", postID))

	_, err := r.Db.Exec(context.Background(), query, updated_post.Title, updated_post.PostContent, updated_post.ReplyTo, postID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("Content of post with id %v successfully updated", postID))
	return nil
}

// Delete one post from database matching post_id. Returns nil if successful.
// Soft delete is performed.
func (r *PostsRepository) Delete(postID string) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)
	utils.Logger.Trace().Msg(fmt.Sprintf("Transaction begin. Deleting post with id: %s", postID))

	// Remove post from posts table
	query := fmt.Sprintf(`
		UPDATE %s SET is_available = false WHERE post_id = $1 AND is_available = true RETURNING author;`, POSTS_TABLE)

	var author string
	if err = tx.QueryRow(ctx, query, postID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting post from database")
		return err
	}

	if author == "" {
		utils.Logger.Warn().Msg("No rows deleted")
		return fmt.Errorf("no rows deleted")
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("Post with id %s successfully deleted from database", postID))

	// Remove post from likes table
	query = fmt.Sprintf(`
	DELETE FROM %s WHERE content_id = $1;`, LIKES_TABLE)
	if _, err = tx.Exec(ctx, query, postID); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting post from likes table")
		return err
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("Post with id %s deleted from likes table", postID))

	// Update user karma
	query = fmt.Sprintf(`
		UPDATE %s SET karma = GREATEST(karma - %d, 0) WHERE uid = $1;`, USERS_TABLE, models.CREATE_POST_PTS)
	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating user karma")
		return err
	}
	utils.Logger.Trace().Msg("User karma successfully updated")

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}
	utils.Logger.Debug().Msg(fmt.Sprintf("Post with id %s successfully deleted from database", postID))
	return nil
}

// Restore deleted post by post_id. Returns nil if successful.
func (r *PostsRepository) Restore(postID string) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)
	utils.Logger.Trace().Msg(fmt.Sprintf("Transaction begin. Restoring post with id: %v", postID))

	// Restore post from posts table
	query := fmt.Sprintf(`
		UPDATE %s SET is_available = true WHERE post_id = $1 AND is_available = false RETURNING author;`, POSTS_TABLE)

	var author string
	if err := tx.QueryRow(ctx, query, postID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring post from database")
		return err
	}
	if author == "" {
		utils.Logger.Warn().Msg("No rows restored")
		return fmt.Errorf("no rows restored")
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("Post with id %s successfully restored from database", postID))

	// Restore user's karma
	query = fmt.Sprintf(`
		UPDATE %s SET karma = karma + %d WHERE uid = $1;`, USERS_TABLE, models.CREATE_POST_PTS)

	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring user karma")
		return err
	}
	utils.Logger.Trace().Msg(fmt.Sprintf("User %s karma successfully restored", author))

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("Post with id %s successfully restored", postID))
	return nil
}

// Restore all deleted posts by thread_id. Returns nil if successful.
func (r *PostsRepository) RestorePostsByThread(threadId string) error {
	query := fmt.Sprintf(`
	WITH restored_rows AS (
		UPDATE %s SET is_available = true WHERE thread_id = $1 AND is_available = false RETURNING thread_id
	)
	SELECT COUNT(*) FROM restored_rows;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Restoring posts with thread_id: %s", threadId))

	var num_restored int
	err := r.Db.QueryRow(context.Background(), query, threadId).Scan(&num_restored)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring posts")
		return err
	}

	utils.Logger.Debug().Int("Rows restored", num_restored).Msg(fmt.Sprintf("Posts with thread_id %s successfully restored", threadId))
	return nil
}
