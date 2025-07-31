package repositories

import (
	"context"
	"errors"
	"fmt"

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
func (r *ThreadsRepository) Insert(thread *models.DbThread) error {
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
	INSERT INTO %s (thread_id, author, title, preview, is_anon)
	VALUES ($1, $2, $3, $4, $5);`, THREADS_TABLE)

	if _, err = tx.Exec(ctx, query, thread.ThreadID, thread.AuthorUid, thread.Title, thread.Preview, thread.IsAnon); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into database")
		return err
	}
	utils.Logger.Debug().Str("thread id", thread.ThreadID).Msg("")

	// Update author's karma
	query = fmt.Sprintf(`
	UPDATE %s
	SET karma = karma + %d
	WHERE uid = $1;`, USERS_TABLE, models.CREATE_THREAD_PTS)
	if _, err = tx.Exec(ctx, query, thread.AuthorUid); err != nil {
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
// uid: user ID of the user requesting the threads; used to determine if the thread is liked by the user
// page: page number - offset is automatically calculated in this function.
// size: page size; number of items to return
// descending: true if sorting is descending, false if ascending
func (r *ThreadsRepository) GetAll(column models.SortColumn, uid string, page int, size int, descending bool) ([]models.Thread, error) {
	desc := "DESC"
	if !descending {
		desc = "ASC"
	}

	// Calculate offset
	offset := (page - 1) * size

	// Query params
	// $1: user's uid
	// $2: limit
	// $3: offset
	query := fmt.Sprintf(`
	SELECT
		T.THREAD_ID,
		T.TITLE,
		-- Conditionally return author UID or 'NA'
		CASE 
			WHEN T.IS_ANON THEN 'NA'
			ELSE T.AUTHOR
		END AS AUTHOR,
		T.TIME_CREATED,
		T.LAST_ACTIVITY,
		(
			SELECT 
				COUNT(1)
			FROM
				views V
			WHERE
				V.CONTENT_ID=T.THREAD_ID
		) AS VIEWS,
		T.FLAGGED,
		T.PREVIEW,
		T.IS_AVAILABLE,
		-- Conditionally return author name or '#ANONYMOUS#'
		CASE 
			WHEN T.IS_ANON THEN '#ANONYMOUS#'
			ELSE U.NAME
		END AS AUTHOR_NAME,
		T.IS_ANON,
		T.AUTHOR=$1 AS IS_AUTHOR, -- uid parameter
		(
			SELECT
				COUNT(1) - 1
			FROM
				posts P
			WHERE
				P.THREAD_ID = T.THREAD_ID
				AND P.IS_AVAILABLE = TRUE
		) AS NUM_REPLIES,
		COUNT(L.CONTENT_ID) AS NUM_LIKES,
		MAX(
			CASE
				WHEN L.UID = $1 THEN 1
				ELSE 0
			END
		)::BOOLEAN AS IS_LIKED,
		MAX(
			CASE
				WHEN F.UID = $1 THEN 1
				ELSE 0
			END
		)::BOOLEAN AS IS_FAVORITED
	FROM
		THREADS T
		INNER JOIN USERS U ON T.AUTHOR = U.UID
		LEFT JOIN LIKES L ON T.THREAD_ID = L.CONTENT_ID
		LEFT JOIN FAVORITES F ON T.THREAD_ID = F.CONTENT_ID
	WHERE
		T.IS_AVAILABLE = TRUE
	GROUP BY
		T.THREAD_ID,
		U.UID
	ORDER BY
		%s %s
	LIMIT $2
	OFFSET $3;`, column, desc)

	utils.Logger.Debug().Str("column", string(column)).Int("page", page).Int("offset", offset).Int("size", size).Bool("descending", descending).Msg("")

	// Perform query and collect rows into array.
	rows, _ := r.Db.Query(context.Background(), query, uid, size, offset)
	threads, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Thread])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Msgf("%d threads retrieved from database", len(threads))

	return threads, nil
}

// Get threads metadata
func (r *ThreadsRepository) GetMetadata() (*models.ContentMetadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) AS COUNT FROM %s WHERE IS_AVAILABLE=TRUE;`, THREADS_TABLE)

	row, _ := r.Db.Query(context.Background(), query)
	defer row.Close()
	metadata, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.ContentMetadata])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Int("num threads", metadata.Total).Msg("Threads metadata retrieved from database")

	return metadata, nil
}

// Get thread and corresponding posts by thread_id. Returns thread object if found, nil otherwise.
func (r *ThreadsRepository) GetByID(thread_id string, uid string) (*models.Thread, error) {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return nil, err
	}

	defer tx.Rollback(ctx)

	// Retrieve the thread
	query := fmt.Sprintf(`
	SELECT
		T.THREAD_ID,
		T.TITLE,
		-- Conditionally return author UID or 'NA'
		CASE 
			WHEN T.IS_ANON THEN 'NA'
			ELSE T.AUTHOR
		END AS AUTHOR,
		T.TIME_CREATED,
		T.LAST_ACTIVITY,
		(
			SELECT 
				COUNT(1) + 1
			FROM
				views V
			WHERE
				V.CONTENT_ID=T.THREAD_ID
		) AS VIEWS,
		T.FLAGGED,
		T.PREVIEW,
		T.IS_AVAILABLE,
		-- Conditionally return author name or '#ANONYMOUS#'
		CASE 
			WHEN T.IS_ANON THEN '#ANONYMOUS#'
			ELSE USERS.NAME
		END AS AUTHOR_NAME,
		T.IS_ANON,
		T.AUTHOR=$1 AS IS_AUTHOR,
		(
			SELECT
				COUNT(1) - 1
			FROM
				POSTS P
			WHERE
				P.THREAD_ID = T.THREAD_ID
				AND P.IS_AVAILABLE = TRUE
		) AS NUM_REPLIES,
		COUNT(L.CONTENT_ID) AS NUM_LIKES,
		MAX(
			CASE
				WHEN L.UID = $2 THEN 1
				ELSE 0
			END
		)::BOOLEAN AS IS_LIKED,
		MAX(
			CASE
				WHEN F.UID = $2 THEN 1
				ELSE 0
			END
		)::BOOLEAN AS IS_FAVORITED
	FROM
		%s T
		INNER JOIN USERS ON T.AUTHOR = USERS.UID
		LEFT JOIN LIKES L ON T.THREAD_ID = L.CONTENT_ID
		LEFT JOIN FAVORITES F ON T.THREAD_ID = F.CONTENT_ID
	WHERE
		T.IS_AVAILABLE = TRUE
		AND T.THREAD_ID = $1
	GROUP BY
		T.THREAD_ID,
		T.TITLE,
		T.AUTHOR,
		T.TIME_CREATED,
		T.LAST_ACTIVITY,
		T.FLAGGED,
		T.PREVIEW,
		T.IS_AVAILABLE,
		T.IS_ANON,
		USERS.NAME;`, THREADS_TABLE)

	utils.Logger.Trace().Msg(fmt.Sprintf("Getting thread with id: %v", thread_id))

	row, _ := tx.Query(context.Background(), query, thread_id, uid)
	thread, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Thread])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting thread by ID")
		return nil, err
	}
	defer row.Close()

	// Insert view into views table
	query = "INSERT INTO VIEWS VALUES ($1, $2) ON CONFLICT DO NOTHING;"
	if _, err = tx.Exec(ctx, query, uid, thread_id); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting into views for thread id %s", thread_id)
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Str("thread_id", thread_id).Msgf("Error committing transaction for thread id %s", thread_id)
		return nil, err
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("Thread with id %v found", thread_id))
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

	if err = tx.QueryRow(ctx, query, threadID).Scan(&author); err != nil {
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
