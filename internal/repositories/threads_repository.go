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

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return thread_id, nil
}

// Get thread author
func (r *ThreadRepository) GetThreadAuthor(thread_id string) (string, error) {
	query := fmt.Sprintf(`SELECT author FROM %s WHERE thread_id = $1;`, THREADS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Getting author of thread with id: %v", thread_id))

	var author string
	err := r.Db.QueryRow(context.Background(), query, thread_id).Scan(&author)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return "", err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Author of thread with id %v is %v", thread_id, author))
	return author, nil
}

// Perform soft delete of the thread in the database. Returns nil if successful.
func (r *ThreadRepository) DeleteThread(thread_id string) error {
	query := fmt.Sprintf(`
	WITH deleted_rows AS (
		UPDATE %s
		SET is_available = false
		WHERE thread_id = $1
		RETURNING thread_id
	)
	SELECT COUNT(*) FROM deleted_rows;`, THREADS_TABLE)

	var num_deleted int
	err := r.Db.QueryRow(context.Background(), query, thread_id).Scan(&num_deleted)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		return err
	} else if num_deleted == 0 {
		utils.Logger.Warn().Msg("No rows affected. Thread may not exist or already deleted")
		return nil
	}

	utils.Logger.Info().Int("num of rows deleted", num_deleted).Msg(fmt.Sprintf("%s successfully deleted", thread_id))
	return nil
}

// (NOT USED)
// Perform hard delete of the thread in the database. Returns nil if successful.
func (r *ThreadRepository) HardDeleteThread(thread_id string) error {
	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE thread_id = $1;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, thread_id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully deleted (HARD DELETE)", thread_id))
	return nil
}

// Restore deleted thread by thread_id. Returns nil if successful.
// Currently unused.
func (r *ThreadRepository) RestoreThread(thread_id string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET is_available = true
	WHERE thread_id = $1;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, thread_id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring thread")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully restored", thread_id))
	return nil
}
