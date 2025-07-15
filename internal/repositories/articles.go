package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Articles table name in db
const ARTICLES_TABLE = "articles"

type ArticleRepository struct {
	_  ContentRepository
	Db *pgxpool.Pool
}

var Articles *ArticleRepository

// Insert new article into database.
func (r *ArticleRepository) Insert(article *models.DbArticle) error {
	ctx := context.Background()

	utils.Logger.Trace().Interface("article", article).Msg("")

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction for inserting article")
		return err
	}
	defer tx.Rollback(ctx)

	// Insert article into articles table
	query := fmt.Sprintf(`
		INSERT INTO %s (article_id, author, title, content, preview)
		VALUES ($1, $2, $3, $4, $5)
	`, ARTICLES_TABLE)

	if _, err := tx.Exec(ctx, query, article.ArticleID, article.Author, article.Title, article.Content, article.Preview); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting article into database")
		return err
	}
	utils.Logger.Trace().Msgf("Inserted article %s into database", article.ArticleID)

	// Update author's karma
	query = fmt.Sprintf(`
	UPDATE %s
	SET karma = karma + %d
	WHERE uid = $1`, USERS_TABLE, models.CREATE_ARTICLE_PTS)

	if _, err := tx.Exec(ctx, query, article.Author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma after article creation")
		return err
	}
	utils.Logger.Trace().Msgf("Updated karma for author %s after article creation", article.Author)

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction after inserting article")
		return err
	}

	utils.Logger.Debug().Msgf("Successfully inserted article %s and updated author's karma", article.ArticleID)
	return nil
}

// GetAll retrieve all articles from a certain timestamp.
func (r *ArticleRepository) GetAll(uid string, column models.ThreadColumn, cursor time.Time, size int, descending bool) ([]models.Article, error) {
	// Set descending string
	desc := "DESC"
	if !descending {
		desc = "ASC"
	}

	query := fmt.Sprintf(`
	SELECT
		A.ARTICLE_ID,
		A.AUTHOR,
		A.TITLE,
		A.TIME_CREATED,
		A.LAST_ACTIVITY,
		A.VIEWS,
		A.FLAGGED,
		A.IS_AVAILABLE,
		A.CONTENT,
		A.PREVIEW,
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
		)::BOOLEAN AS IS_LIKED
	FROM
		ARTICLES A
		INNER JOIN USERS U ON A.AUTHOR = U.UID
		LEFT JOIN LIKES L ON A.ARTICLE_ID = L.CONTENT_ID
	WHERE
		A.%s < $2		-- Cursor parameter
		AND A.IS_AVAILABLE = TRUE
	GROUP BY
		A.ARTICLE_ID,
		U.UID
	ORDER BY
		A.%s %s		-- Column, ASC/DESC
	LIMIT
		$3;		-- Size parameter
	`, column, column, desc)

	utils.Logger.Debug().Str("column", string(column)).Time("cursor", cursor).Int("size", size).Bool("descenting", descending).Msg("article query params")

	rows, _ := r.Db.Query(context.Background(), query, uid, cursor, size)
	articles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Article])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Msgf("%d articles retrieved from database", len(articles))

	return articles, nil
}

// Retrieve an article by its ID and increment its view count.
// It also returns the article's author name, number of comments, number of likes, and
// whether the article is liked by the user with the given UID.
// If the article is not found or an error occurs, it returns an error.
func (r *ArticleRepository) GetByID(articleID string, uid string) (*models.Article, error) {
	ctx := context.Background()

	utils.Logger.Trace().Msgf("Fetching article with ID %s for user %s", articleID, uid)

	query := fmt.Sprintf(`
	WITH
	A AS (
		UPDATE %s
		SET
			VIEWS = VIEWS + 1
		WHERE
			ARTICLE_ID = $1
			AND IS_AVAILABLE = TRUE
		RETURNING
			*
	)
	SELECT
		A.ARTICLE_ID,
		A.AUTHOR,
		A.TITLE,
		A.TIME_CREATED,
		A.LAST_ACTIVITY,
		A.VIEWS,
		A.FLAGGED,
		A.IS_AVAILABLE,
		A.CONTENT,
		A.PREVIEW,
		USERS.NAME AS AUTHOR_NAME,
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
				WHEN L.UID = $2 THEN 1
				ELSE 0
			END
		)::BOOLEAN AS IS_LIKED
	FROM
		A
		INNER JOIN USERS ON A.AUTHOR = USERS.UID
		LEFT JOIN LIKES L ON A.ARTICLE_ID = L.CONTENT_ID
	WHERE
		A.IS_AVAILABLE = TRUE
		AND A.ARTICLE_ID = $1
	GROUP BY
		A.ARTICLE_ID,
		A.AUTHOR,
		A.TITLE,
		A.TIME_CREATED,
		A.LAST_ACTIVITY,
		A.VIEWS,
		A.FLAGGED,
		A.IS_AVAILABLE,
		A.CONTENT,
		A.PREVIEW,
		USERS.NAME;
	`, ARTICLES_TABLE)

	utils.Logger.Trace().Msgf("Getting article with id: %s", articleID)

	row, _ := r.Db.Query(ctx, query, articleID, uid)
	defer row.Close()
	article, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Article])
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.Logger.Warn().Msgf("No article found with ID %s", articleID)
		}

		utils.Logger.Error().Err(err).Msgf("Error fetching article with ID %s",
			articleID)
		return nil, err
	}

	utils.Logger.Debug().Msgf("Successfully fetched article with ID %s", articleID)
	return article, nil

}

// Get articles metadata
func (r *ArticleRepository) GetMetadata() (*models.ArticlesMetadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) AS NUM_ARTICLES FROM %s WHERE IS_AVAILABLE=TRUE;`, ARTICLES_TABLE)

	row, _ := r.Db.Query(context.Background(), query)
	defer row.Close()
	metadata, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.ArticlesMetadata])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error collecting rows")
		return nil, err
	}

	utils.Logger.Debug().Int("num articles", metadata.NumArticles).Msg("Article metadata retrieved from database")

	return metadata, nil
}

// Get article author by article ID.
// Returns the author's UID if found, otherwise empty string.
func (r *ArticleRepository) GetAuthor(articleID string) (string, error) {
	ctx := context.Background()

	utils.Logger.Trace().Msgf("Fetching author for article with ID %s", articleID)

	query := fmt.Sprintf(`
	SELECT AUTHOR
	FROM %s
	WHERE ARTICLE_ID = $1 AND IS_AVAILABLE = TRUE;
	`, ARTICLES_TABLE)

	var author string
	if err := r.Db.QueryRow(ctx, query, articleID).Scan(&author); err != nil {
		if err == pgx.ErrNoRows {
			utils.Logger.Warn().Msgf("No author found for article with ID %s", articleID)
			return "", nil
		}
		utils.Logger.Error().Err(err).Msgf("Error fetching author for article with ID %s", articleID)
		return "", err
	}

	utils.Logger.Info().Msgf("Successfully fetched author %s for article with ID %s", author, articleID)
	return author, nil
}

// Returns true if article exists in database
func (r *ArticleRepository) IsAvailable(articleID string) bool {
	query := fmt.Sprintf(`SELECT IS_AVAILABLE FROM %s WHERE ARTICLE_ID=$1;`, ARTICLES_TABLE)

	var isAvailable bool
	if err := r.Db.QueryRow(context.Background(), query, articleID).Scan(&isAvailable); err != nil {
		utils.Logger.Warn().Msgf("Article of ID %s not available", articleID)
		return false
	}

	return isAvailable
}

// Perform soft delete of an article by its ID.
// It also updates the author's karma by subtracting the points for article creation.
func (r *ArticleRepository) Delete(articleID string) error {
	ctx := context.Background()

	utils.Logger.Trace().Msgf("Deleting article with ID %s", articleID)

	// Begin transaction
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction for deleting article")
		return err
	}
	defer tx.Rollback(ctx)

	// Author's uid
	var author string
	// Soft delete the article
	query := fmt.Sprintf(`
		UPDATE %s
		SET IS_AVAILABLE = FALSE, last_activity = NOW()
		WHERE ARTICLE_ID = $1
		AND IS_AVAILABLE = TRUE
		RETURNING AUTHOR;
	`, ARTICLES_TABLE)

	if err := tx.QueryRow(ctx, query, articleID).Scan(&author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error soft deleting article")
		return err
	}

	utils.Logger.Trace().Msgf("Soft deleted article with ID %s", articleID)

	// Update author's karma
	query = fmt.Sprintf(`
	UPDATE %s
	SET karma = GREATEST(karma - %d, 0)
	WHERE uid = $1`, USERS_TABLE, models.CREATE_ARTICLE_PTS)

	if _, err := tx.Exec(ctx, query, author); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating author's karma after article deletion")
		return err
	}

	utils.Logger.Trace().Msgf("Updated karma for author %s after article deletion", author)

	// Remove article from likes table
	query = fmt.Sprintf(`
	DELETE FROM %s
	WHERE CONTENT_ID = $1;`, LIKES_TABLE)

	if _, err := tx.Exec(ctx, query, articleID); err != nil {
		utils.Logger.Error().Err(err).Msg("Error removing article from likes table")
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction after deleting article")
		return err
	}

	utils.Logger.Debug().Msgf("Successfully deleted article with ID %s and updated author's karma", articleID)
	return nil
}
