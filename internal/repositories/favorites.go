package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const FAVORITES_TABLE = "favorites"

type FavoritesRepository struct {
	Db *pgxpool.Pool
}

var Favorites *FavoritesRepository

// Insert new favorite into the database, return nil on successful insert.
func (r *FavoritesRepository) Insert(favorite *models.Favorite) error {
	ctx := context.Background()

	query := fmt.Sprintf(`INSERT INTO %s (UID, CONTENT_ID) VALUES ($1, $2) ON CONFLICT (UID, CONTENT_ID) DO NOTHING;`, FAVORITES_TABLE)

	if _, err := r.Db.Exec(ctx, query, favorite.Uid, favorite.ContentId); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting favorite into database for %s", favorite.Uid)
		return err
	}

	utils.Logger.Debug().Str("uid", favorite.Uid).Str("content id", favorite.ContentId).Msgf("Successfully inserted into favorites table for %s", favorite.Uid)

	return nil
}

// Get favorite by uid and content_id, return true if record exists in database and false otherwise.
func (r *FavoritesRepository) Exists(uid string, contentID string) bool {
	query := fmt.Sprintf(`SELECT 1 FROM %s WHERE uid = $1 AND content_id = $2;`, FAVORITES_TABLE)

	var numRecords int
	err := r.Db.QueryRow(context.Background(), query, uid, contentID).Scan(&numRecords)
	if numRecords == 0 || err != nil {
		utils.Logger.Trace().Str("uid", uid).Str("content_id", contentID).Msg("Record not found")
		return false
	}

	utils.Logger.Trace().Str("uid", uid).Str("content_id", contentID).Msg("Record found")
	return true
}

// Get saved threads
func (r *FavoritesRepository) GetThreads(uid string) ([]models.Thread, error) {
	query := `
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
		-- Conditionally return author name or 'ANONYMOUS'
		CASE 
			WHEN T.IS_ANON THEN 'ANONYMOUS'
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
		AND F.UID = $1
	GROUP BY
		T.THREAD_ID,
		U.UID,
		F.TIMESTAMP
	ORDER BY
		F.TIMESTAMP DESC;`

	rows, _ := r.Db.Query(context.Background(), query, uid)
	threads, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Thread])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Msgf("%d saved threads retrieved from database.", len(threads))

	return threads, nil
}

// Get saved articles
func (r *FavoritesRepository) GetArticles(uid string) ([]models.Article, error) {
	query := `
	SELECT
		A.ARTICLE_ID,
		A.AUTHOR,
		A.TITLE,
		A.TIME_CREATED,
		A.LAST_ACTIVITY,
		(
			SELECT 
				COUNT(1)
			FROM
				views V
			WHERE
				V.CONTENT_ID=A.ARTICLE_ID
		) AS VIEWS,
		A.FLAGGED,
		A.IS_AVAILABLE,
		A.CONTENT,
		A.PREVIEW,
		A.AUTHOR=$1 AS IS_AUTHOR, -- uid parameter
		U.NAME AUTHOR_NAME,
		(
			SELECT
				COUNT(1)
			FROM
				COMMENTS C
			WHERE
				C.ARTICLE_ID = A.ARTICLE_ID
				AND C.IS_AVAILABLE = TRUE
		) AS NUM_COMMENTS,
		COUNT(L.CONTENT_ID) AS NUM_LIKES,
		MAX(
			CASE
				WHEN L.UID = $1 THEN 1 	-- User UID parameter
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
		ARTICLES A
		INNER JOIN USERS U ON A.AUTHOR = U.UID
		LEFT JOIN LIKES L ON A.ARTICLE_ID = L.CONTENT_ID
		LEFT JOIN FAVORITES F ON A.ARTICLE_ID = F.CONTENT_ID
	WHERE
		A.IS_AVAILABLE = TRUE
		AND F.UID = $1
	GROUP BY
		A.ARTICLE_ID,
		U.UID,
		F.TIMESTAMP
	ORDER BY 
		F.TIMESTAMP DESC;`

	rows, _ := r.Db.Query(context.Background(), query, uid)
	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Article])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Msgf("%d saved articles retrieved from database.", len(articles))

	return articles, nil
}

// Remove favorite from database (hard delete is performed)
// Returns nil on success
func (r *FavoritesRepository) Delete(uid string, contentID string) error {
	ctx := context.Background()

	query := fmt.Sprintf(`DELETE FROM %s WHERE UID=$1 AND CONTENT_ID=$2;`, FAVORITES_TABLE)

	if _, err := r.Db.Exec(ctx, query, uid, contentID); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting favorite into database for %s", uid)
		return err
	}

	utils.Logger.Debug().Str("uid", uid).Str("content id", contentID).Msgf("Successfully removed from favorites table for %s", uid)

	return nil
}
