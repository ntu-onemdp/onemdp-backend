package repositories

import (
	"context"
	"fmt"

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
		INSERT INTO %s (article_id, author, title, content)
		VALUES ($1, $2, $3, $4)
	`, ARTICLES_TABLE)

	if _, err := tx.Exec(ctx, query, article.ArticleID, article.Author, article.Title, article.Content); err != nil {
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
func (r *ArticleRepository) GetAll() {
	panic("not implemented")
}

// Get articles metadata
func (r *ArticleRepository) GetMetadata() (models.ArticlesMetadata, error) {
	panic("not implemented")
}

func (r *ArticleRepository) GetByID(articleID string) (*models.DbArticle, error) {
	panic("not implemented")
}
