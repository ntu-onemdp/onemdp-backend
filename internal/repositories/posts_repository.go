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

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", post.Title))
	return nil
}

// Delete post from the database. Returns nil on successful delete
// Current implementation: Hard delete. Change to soft delete in the future.
func (r *PostsRepository) DeletePost(postId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE post_id = $1;`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Deleting post with id: %v", postId))

	_, err := r.Db.Exec(context.Background(), query, postId)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg(fmt.Sprintf("Post with id %v successfully deleted from database", postId))
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

	utils.Logger.Trace().Msg(fmt.Sprintf("Author of post with id %v is %v", postId, author))
	return author, nil
}
