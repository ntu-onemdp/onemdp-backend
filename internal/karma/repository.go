package karma

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/semester"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const KARMA_TABLE = "karma"

type karmaRepository struct {
	db *pgxpool.Pool
}

// Insert new semester. Call this function whenever a new semester is created. The database will populate it with the default values of 0.
func (r *karmaRepository) insert(semester string) error {
	query := fmt.Sprintf(`INSERT INTO %s (SEMESTER) VALUES ($1);`, KARMA_TABLE)

	if _, err := r.db.Exec(context.Background(), query, semester); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting new semester to karma table")
		return err
	}

	utils.Logger.Info().Str("semester", semester).Msg("New semester inserted into karma table")
	return nil
}

// Get karma settings for current semester. Join operation is performed with SEMESTERS table to obtain current sem automatically
func (r *karmaRepository) getSettings() (*Karma, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE SEMESTER = (SELECT SEMESTER FROM %s WHERE IS_CURRENT);`, KARMA_TABLE, semester.SEMESTER_TABLE)

	row, _ := r.db.Query(context.Background(), query)
	settings, err := pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByName[Karma])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error retrieving karma settings for current semester")
		return nil, err
	}

	utils.Logger.Debug().Str("semester", settings.Semester).Msg("Karma settings for current semester retrieved")
	return settings, nil
}

// Update karma settings for current semester
func (r *karmaRepository) update(settings Karma) error {
	query := fmt.Sprintf(`UPDATE %s SET CREATE_THREAD=$1, CREATE_ARTICLE=$2, CREATE_COMMENT=$3, CREATE_POST=$4, RECEIVE_LIKE=$5 WHERE SEMESTER = (SELECT SEMESTER FROM %s WHERE IS_CURRENT);`, KARMA_TABLE, semester.SEMESTER_TABLE)

	if _, err := r.db.Exec(context.Background(), query, settings.CreateThreadPts, settings.CreateArticlePts, settings.CreateCommentPts, settings.CreatePostPts, settings.LikePts); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating karma settings for current semester")
		return err
	}

	utils.Logger.Info().Str("semester", settings.Semester).Msg("Karma settings for current semester updated")
	return nil
}
