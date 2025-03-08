package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadRepository struct {
	Db *pgxpool.Pool
}

// Threads table name in db
const THREADS_TABLE = "threads"

// Insert new thread into the database. Returns UUID of new thread on successful insert
func (r *ThreadRepository) CreateThread(thread *models.NewThread) (string, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (thread_id, author, title, preview) 
	VALUES ($1, $2, $3, $4);`, THREADS_TABLE)

	// Generate unique ID for the thread.
	thread_id, err := gonanoid.New(6)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error generating unique ID")
		return "", err
	}

	_, err = r.Db.Exec(context.Background(), query, thread_id, thread.Author, thread.Title, thread.Preview)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return "", err
	}
	utils.Logger.Debug().Str("thread id", thread_id).Msg("")

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return thread_id, nil
}
