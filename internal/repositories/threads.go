package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Threads table name in db
const THREADS_TABLE = "threads"

type ThreadsRepository struct {
	_  ContentRepository
	Db *pgxpool.Pool
}

var Threads *ThreadsRepository

// Insert new thread into the database. Returns thread ID and UUID of header post on successful insert
// Although function takes in a thread object, only author, title and preview are used.
func (r *ThreadsRepository) Create(thread *models.Thread) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Insert thread into database
	query := fmt.Sprintf(`
	INSERT INTO %s (thread_id, author, title, preview) 
	VALUES ($1, $2, $3, $4);`, THREADS_TABLE)

	if _, err = tx.Exec(ctx, query, thread.ThreadID, thread.Author, thread.Title, thread.Preview); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return err
	}
	utils.Logger.Debug().Str("thread id", thread.ThreadID).Msg("")

	// Update author's karma
	query = fmt.Sprintf(`
	UPDATE %s
	SET karma = karma + %d
	WHERE uid = $1;`, USERS_TABLE, models.CREATE_THREAD_PTS)
	if _, err = tx.Exec(ctx, query, thread.Author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma")
		return err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("%s successfully inserted into database", thread.Title))
	return nil
}

// GetAll all threads from a certain timestamp. Cursor given is the timestamp of the last thread in the previous page.
// Number of threads returned is not hardcoded, can be chosen in frontend.
// Params
// column: column to sort by
// cursor: timestamp of the last thread in the previous page
// size: page size; number of threads to return
// descending: true if sorting is descending, false if ascending
func (r *ThreadsRepository) GetAll(column models.ThreadColumn, cursor time.Time, size int, descending bool) ([]models.Thread, error) {
	desc := "DESC"
	if !descending {
		desc = "ASC"
	}

	// Example SQL statement after formatting:
	// SELECT * FROM threads WHERE time_created < cursor AND is_available = true ORDER BY time_created DESC LIMIT size;
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s < $1 AND is_available = true ORDER BY %s %s LIMIT $2;`, THREADS_TABLE, column, column, desc)

	utils.Logger.Debug().Str("column", string(column)).Time("cursor", cursor).Int("size", size).Bool("descending", descending).Msg("")
	rows, _ := r.Db.Query(context.Background(), query, cursor, size)
	threads, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Thread])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Info().Msg("Threads retrieved from database")
	return threads, nil
}

// Get threads metadata
func (r *ThreadsRepository) GetMetadata() (models.ThreadsMetadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE is_available = true;`, THREADS_TABLE)

	var num_threads int
	err := r.Db.QueryRow(context.Background(), query).Scan(&num_threads)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting threads metadata")
		return models.ThreadsMetadata{}, err
	}

	utils.Logger.Info().Int("num_threads", num_threads).Msg("Threads metadata retrieved")
	return models.ThreadsMetadata{NumThreads: num_threads}, nil
}

// Get thread by thread_id. Returns thread object if found, nil otherwise.
func (r *ThreadsRepository) GetByID(thread_id string) (*models.Thread, error) {
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
func (r *ThreadsRepository) GetAuthor(thread_id string) (string, error) {
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
func (r *ThreadsRepository) IsAvailable(thread_id string) bool {
	query := fmt.Sprintf(`SELECT is_available FROM %s WHERE thread_id = $1;`, THREADS_TABLE)

	var is_available bool
	err := r.Db.QueryRow(context.Background(), query, thread_id).Scan(&is_available)
	if err != nil {
		return false
	}

	return is_available
}

// Update thread's last activity timestamp to current time. Returns nil if successful.
func (r *ThreadsRepository) UpdateActivity(thread_id string) error {
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
func (r *ThreadsRepository) Update(threadID string, title string, preview string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET title = $1, preview = $2, last_activity = NOW()
	WHERE thread_id = $3 AND is_available = true;`, THREADS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, title, preview, threadID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating preview")
		return err
	}

	utils.Logger.Info().Msg(fmt.Sprintf("Thread %s preview updated", threadID))
	return nil
}

// Perform soft delete of the thread in the database. Returns nil if successful.
// Karma rollback:
//  1. Author's karma is decremented by CREATE_THREAD_PTS
//  2. All posts with matching thread_id are matched with their authors and decremented by CREATE_POST_PTS
//  3. We do not decrement the karma of people who had liked the posts in the thread. If you had given a
//     good answer to a post, your karma should be deserved and thus not decremented.
func (r *ThreadsRepository) Delete(threadID string) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	var author string
	query := fmt.Sprintf(`
		UPDATE %s
		SET is_available = false, last_activity = NOW()
		WHERE thread_id = $1 AND is_available = true
		RETURNING author;`, THREADS_TABLE)

	if err = tx.QueryRow(context.Background(), query, threadID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread")
		return errors.New("thread is not available or does not exist")
	} else if author == "" {
		utils.Logger.Warn().Msg("No rows affected. Thread may not exist or already deleted")
		return nil
	}
	utils.Logger.Trace().Msg("Thread deleted")

	// Update author's karma
	query = fmt.Sprintf(`
	UPDATE %s
	SET karma = GREATEST(karma - %d, 0)
	WHERE uid = $1;`, USERS_TABLE, models.CREATE_THREAD_PTS)

	if _, err = tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma")
		return err
	}
	utils.Logger.Trace().Str("author", author).Msg("Karma updated for author of original thread.")

	// Remove thread from likes table
	query = fmt.Sprintf(`
	DELETE FROM %s
	WHERE content_id = $1;`, LIKES_TABLE)

	if _, err = tx.Exec(ctx, query, threadID); err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting thread from likes table")
		return err
	}
	utils.Logger.Trace().Str("thread id", threadID).Msg("Thread deleted from likes table")

	// Remove thread from posts table.
	query = fmt.Sprintf(`
	UPDATE %s
	SET is_available = false, last_edited = NOW()
	WHERE thread_id = $1 AND is_available = true
	RETURNING post_id, author, is_header;`, POSTS_TABLE)

	rows, err := tx.Query(ctx, query, threadID)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error deleting from posts table posts with thread id of " + threadID)
		return err
	}
	defer rows.Close()

	// Prepare batch delete
	batch := &pgx.Batch{}
	for rows.Next() {
		var postID string // Post ID of the post to be deleted
		var author string // Author of the post to be deleted
		var isHeader bool // Whether the post is a header post

		if err = rows.Scan(&postID, &author, &isHeader); err != nil {
			utils.Logger.Error().Err(err).Msg("Error scanning rows")
			return err
		}
		utils.Logger.Trace().Str("post id", postID).Str("author", author).Bool("is header", isHeader).Msg("deleting post")

		// Delete post from likes table
		query = fmt.Sprintf(`
		DELETE FROM %s
		WHERE content_id = $1
		RETURNING 1;
		`, LIKES_TABLE)

		batch.Queue(query, postID)
		utils.Logger.Trace().Str("query", query).Msg("Query added to batch")

		// Update karma of the author of the post
		// Do not update karma of the author of the header post
		if !isHeader {
			query = fmt.Sprintf(`
			UPDATE %s
			SET karma = GREATEST(karma - %d, 0)
			WHERE uid = $1;`, USERS_TABLE, models.CREATE_POST_PTS)

			batch.Queue(query, author)
			utils.Logger.Trace().Str("query", query).Msg("Query added to batch")
		}
	}
	rows.Close()

	// Execute batch delete
	br := tx.SendBatch(ctx, batch)
	br.Close()
	utils.Logger.Trace().Msg("Batch delete executed")

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}
	utils.Logger.Trace().Msg("Transaction committed")

	utils.Logger.Debug().Msg(fmt.Sprintf("%s successfully deleted", threadID))
	return nil
}

// (NOT USED)
// Perform hard delete of the thread in the database. Returns nil if successful.
// func (r *ThreadsRepository) HardDelete(threadID string) error {
// 	query := fmt.Sprintf(`
// 	DELETE FROM %s
// 	WHERE thread_id = $1;`, THREADS_TABLE)

// 	_, err := r.Db.Exec(context.Background(), query, threadID)
// 	if err != nil {
// 		utils.Logger.Error().Err(err).Msg("Error deleting thread")
// 		return err
// 	}

// 	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully deleted (HARD DELETE)", threadID))
// 	return nil
// }

// Restore deleted thread by thread_id. Returns nil if successful.
// Currently unused.
// func (r *ThreadsRepository) Restore(threadID string) error {
// 	query := fmt.Sprintf(`
// 	UPDATE %s
// 	SET is_available = true
// 	WHERE thread_id = $1;`, THREADS_TABLE)

// 	_, err := r.Db.Exec(context.Background(), query, threadID)
// 	if err != nil {
// 		utils.Logger.Error().Err(err).Msg("Error restoring thread")
// 		return err
// 	}

// 	utils.Logger.Info().Msg(fmt.Sprintf("%s successfully restored", threadID))
// 	return nil
// }
