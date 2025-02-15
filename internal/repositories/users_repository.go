package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type UsersRepository struct {
	Db *pgxpool.Pool
}

// Retrieve user's status based on the username.
// Throws error if username cannot be found
func (r *UsersRepository) GetStatusByUsername(username string) (string, error) {
	query := "SELECT status FROM users WHERE username=$1;"

	var status string
	err := r.Db.QueryRow(context.Background(), query, username).Scan(&status)
	if err != nil {
		utils.Logger.Error().Err(err)
		return "", err
	}

	return status, nil
}

// Retrieve user's information from username
// This function *checks* if user is active before returning. If the user's status is not 'active',
// an error is return instead.
func (r *UsersRepository) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username=$1 AND status = 'active';"

	row, _ := r.Db.Query(context.Background(), query, username)
	user, err := pgx.RowToStructByName[models.User](row)

	if err != nil {
		utils.Logger.Error().Err(err)
		return nil, err
	}

	return &user, nil
}

// This method does not work for now. Explore in the future when there is time.
// func (r *UsersRepository) InsertManyUsers(users []models.User) error {

// 	// Open a transaction
// 	tx, err := r.Db.Begin(context.Background())
// 	if err != nil {
// 		utils.Logger.Error().Err(err)
// 		return err
// 	}
// 	defer tx.Rollback(context.Background())

// 	utils.Logger.Trace().Msg("fn called")
// 	// Create temporary staging table
// 	STAGING_USERS := "staging_users" // temp staging table name
// 	_, err = r.Db.Exec(context.Background(), "CREATE TEMP TABLE staging_users (username TEXT PRIMARY KEY,name TEXT);")
// 	if err != nil {
// 		utils.Logger.Error().Err(err)
// 		return err
// 	}

// 	// Prepare copy operation
// 	copyCount, err := tx.CopyFrom(context.Background(), pgx.Identifier{STAGING_USERS}, []string{"username, name"}, pgx.CopyFromSlice(len(users), func(i int) ([]any, error) {
// 		return []any{users[i].Username, users[i].Name}, nil
// 	}))

// 	if err != nil {
// 		utils.Logger.Error().Err(err)
// 		return err
// 	}

// 	utils.Logger.Debug().Msg("Users copied to staging table")

// 	// Move data to final table
// 	_, err = tx.Exec(context.Background(), `
// 	INSERT INTO users (username, name)
// 	SELECT s.username, s.name
// 	FROM staging_users s
// 	WHERE NOT EXISTS (
// 		SELECT 1 FROM users u WHERE u.username = s.username
// 		)
// 		ON CONFLICT (username) DO UPDATE
// 		SET name = EXCLUDED.name
// 		`)

// 	if err != nil {
// 		utils.Logger.Error().Err(err)
// 		return err
// 	}

// 	utils.Logger.Trace().Msg(strconv.FormatInt(copyCount, 10))

// 	return tx.Commit(context.Background())
// }
