package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadRepository struct {
	Db *pgxpool.Pool
}

// Threads table name in db
const THREADS_TABLE = "threads"

// Insert new thread into the database. Returns nil on successful insert
func (r *ThreadRepository) CreateThread(thread *models.NewThread) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (author, title, preview) 
	VALUES ($1, $2, $3);`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, thread.Author, thread.Title, thread.Preview)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return nil
}
