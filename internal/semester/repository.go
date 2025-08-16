package semester

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

const SEMESTER_TABLE = "semesters"

type SemesterRepository struct {
	db *pgxpool.Pool
}

var repo *SemesterRepository

// Insert new sem
func (r *SemesterRepository) insert(semester Semester) error {
	ctx := context.Background()

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Set all is_current of ALL semesters to false
	query := fmt.Sprintf(`UPDATE %s SET IS_CURRENT=FALSE;`, SEMESTER_TABLE)

	if _, err := tx.Exec(ctx, query); err != nil {
		utils.Logger.Error().Err(err).Msg("Error setting is_current of all semesters to false")
		return err
	}

	// Insert new semester into database
	query = fmt.Sprintf(`INSERT INTO %s (SEMESTER, CODE, IS_CURRENT) VALUES $1, $2, $3;`, SEMESTER_TABLE)
	if _, err := tx.Exec(ctx, query, semester.Semester, semester.Code, semester.IsCurrent); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting new semester into database")
		return err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Error committing transaction")
		return err
	}

	utils.Logger.Info().Str("Semester", semester.Semester).Str("Code", semester.Code).Msgf("New semester %s successfully updated in database.", semester.Semester)
	return nil
}

// Retreive current sem
func (r *SemesterRepository) getCurrentSem() (*Semester, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE IS_CURRENT;`, SEMESTER_TABLE)

	row, _ := r.db.Query(context.Background(), query)
	semester, err := pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByName[Semester])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error retrieving code for the current semester")
		return nil, err
	}

	utils.Logger.Debug().Str("code", semester.Code).Msg("Code for current semester retrieved from db.")
	return semester, nil
}

// Set new code for current sem
func (r *SemesterRepository) RefreshCode(code string) error {
	query := fmt.Sprintf(`UPDATE %s SET CODE=$1 WHERE IS_CURRENT;`, SEMESTER_TABLE)

	res, err := r.db.Exec(context.Background(), query, code)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating code for current semester")
		return err
	}

	if res.RowsAffected() != 1 {
		utils.Logger.Error().Int64("Rows affected", res.RowsAffected()).Msg("Number of rows updated is not 1")

		if res.RowsAffected() > 1 {
			return pgx.ErrTooManyRows
		} else {
			return pgx.ErrNoRows
		}
	}

	utils.Logger.Info().Str("New code", code).Msg("Code successfully updated for current semester")
	return nil
}
