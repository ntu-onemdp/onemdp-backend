package repositories

import (
	"context"
	"fmt"

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
	INSERT INTO %s (author, thread_id, title, content) 
	VALUES ($1, $2, $3, $4);`, POSTS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Inserting post: %v", post))

	_, err := r.Db.Exec(context.Background(), query, post.Author, post.ThreadId, post.Title, post.Content)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", post.Title))
	return nil
}
