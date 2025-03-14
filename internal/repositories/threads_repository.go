package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadRepository struct {
	ContentRepository
	Db *pgxpool.Pool
}

// Threads table name in db
const THREADS_TABLE = "threads"

// Insert new thread into the database. Returns thread ID and UUID of header post on successful insert
// Although function takes in a thread object, only author, title and preview are used.
func (r *ThreadRepository) Create(thread *models.Thread) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (thread_id, author, title, preview) 
	VALUES ($1, $2, $3, $4);`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, thread.GetID(), thread.GetAuthor(), thread.GetTitle(), thread.Preview)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return err
	}
	utils.Logger.Debug().Str("thread id", thread.ThreadID).Msg("")

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return nil
}

// Get thread by thread_id. Returns thread object if found, nil otherwise.
func (r *ThreadRepository) GetByID(thread_id string) (*models.Thread, error) {
	// This function is called only when a thread is requested by its ID, so views are incremented here
	query := fmt.Sprintf(`
	WITH thread AS (
		UPDATE %s
		SET views = views + 1
		WHERE thread_id = $1 AND is_available = true
		RETURNING *
	)
	SELECT * FROM thread;`, THREADS_TABLE)

	utils.Logger.Debug().Msg(fmt.Sprintf("Getting thread with id: %v", thread_id))

	row, _ := r.Db.Query(context.Background(), query, thread_id)
	thread, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Thread])
	if err != nil {
		return nil, err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Thread with id %v found", thread_id))
	return &thread, nil
}

// Get thread author
func (r *ThreadRepository) GetAuthor(thread_id string) (string, error) {
	query := fmt.Sprintf(`SELECT author FROM %s WHERE thread_id = $1 AND is_available = true;`, THREADS_TABLE)

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

// Returns true if the thread exists
func (r *ThreadRepository) IsAvailable(thread_id string) bool {
	query := fmt.Sprintf(`SELECT is_available FROM %s WHERE thread_id = $1;`, THREADS_TABLE)

	var is_available bool
	err := r.Db.QueryRow(context.Background(), query, thread_id).Scan(&is_available)
	if err != nil {
		return false
	}

	return is_available
}

// Update thread's last activity timestamp to current time. Returns nil if successful.
func (r *ThreadRepository) UpdateActivity(thread_id string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET last_activity = NOW()
	WHERE thread_id = $1;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, thread_id)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating last activity")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Thread %s last activity updated", thread_id))
	return nil
}

// Update thread's title and preview. Returns nil if successful.
// Thread's title and preview is updated only when the header post is updated.
func (r *ThreadRepository) Update(threadID string, title string, preview string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET title = $1, preview = $2, last_activity = NOW()
	WHERE thread_id = $3;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, title, preview, threadID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating preview")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Thread %s preview updated", threadID))
	return nil
}

// Perform soft delete of the thread in the database. Returns nil if successful.
func (r *ThreadRepository) Delete(threadID string) error {
	query := fmt.Sprintf(`
	WITH deleted_rows AS (
		UPDATE %s
		SET is_available = false, last_activity = NOW()
		WHERE thread_id = $1
		RETURNING thread_id
	)
	SELECT COUNT(*) FROM deleted_rows;`, THREADS_TABLE)

	var num_deleted int
	err := r.Db.QueryRow(context.Background(), query, threadID).Scan(&num_deleted)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		return err
	} else if num_deleted == 0 {
		utils.Logger.Warn().Msg("No rows affected. Thread may not exist or already deleted")
		return nil
	}

	utils.Logger.Info().Int("num of rows deleted", num_deleted).Msg(fmt.Sprintf("%s successfully deleted", threadID))
	return nil
}

// (NOT USED)
// Perform hard delete of the thread in the database. Returns nil if successful.
func (r *ThreadRepository) HardDelete(threadID string) error {
	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE thread_id = $1;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, threadID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully deleted (HARD DELETE)", threadID))
	return nil
}

// Restore deleted thread by thread_id. Returns nil if successful.
// Currently unused.
func (r *ThreadRepository) Restore(threadID string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET is_available = true
	WHERE thread_id = $1;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, threadID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error restoring thread")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully restored", threadID))
	return nil
}
