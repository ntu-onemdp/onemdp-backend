package repositories

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadRepository struct {
	Db *pgxpool.Pool
}

// Threads table name in db
const THREADS_TABLE = "threads"

// Insert new thread into the database. Returns UUID of new thread on successful insert
func (r *ThreadRepository) CreateThread(thread *models.NewThread) (uuid.UUID, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (author, title, preview) 
	VALUES ($1, $2, $3)
	RETURNING "thread_id";`, THREADS_TABLE)

	var threadId uuid.UUID
	err := r.Db.QueryRow(context.Background(), query, thread.Author, thread.Title, thread.Preview).Scan(&threadId)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return uuid.Nil, err
	}
	utils.Logger.Debug().Msg(fmt.Sprintf("ThreadId: %s", threadId.String()))

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return threadId, nil
}
