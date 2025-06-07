package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Users table name in db
const USERS_TABLE = "users"

type UsersRepository struct {
	Db *pgxpool.Pool
}

var Users *UsersRepository

// Insert one empty user into the database. Returns nil on successful insert
// Use this function for user creation
func (r *UsersRepository) InsertOneUser(user *models.User) error {
	query := `
	INSERT INTO users (email, role, semester) 
	VALUES ($1, $2, $3);`

	_, err := r.Db.Exec(context.Background(), query, user.Email, user.Role, user.Semester)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", user.Email))
	return nil
}

// Admin: Retrieve all users from the database, regardless of status
func (r *UsersRepository) GetAllUsers() ([]models.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s;`, USERS_TABLE)
	rows, err := r.Db.Query(context.Background(), query)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error serializing rows to user structs")
		return nil, err
	}

	return users, nil
}

// Retrieve user's status based on the email.
// Throws error if email cannot be found
func (r *UsersRepository) GetStatusByEmail(email string) (string, error) {
	query := fmt.Sprintf(`SELECT status FROM %s WHERE email=$1;`, USERS_TABLE)

	var status string
	err := r.Db.QueryRow(context.Background(), query, email).Scan(&status)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return "", err
	}

	return status, nil
}

// Retrieve user's information from username
// This function *checks* if user is active before returning. If the user's status is not 'active',
// an error is return instead.
func (r *UsersRepository) GetUserByUsername(username string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND status='active';", USERS_TABLE)
	row, _ := r.Db.Query(context.Background(), query, username)
	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	utils.Logger.Debug().Str("username", username).Msg("")
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}
	return user, nil
}

// Admin: Retrieve user's information from username
// This function is able to retrieve deleted users as well
func (r *UsersRepository) GetUserByUsernameAdmin(username string) (*models.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1;`, USERS_TABLE)
	row, _ := r.Db.Query(context.Background(), query, username)
	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	return user, nil
}

// Retrieve if user has changed password
func (r *UsersRepository) GetUserPasswordChanged(username string) (bool, error) {
	query := fmt.Sprintf(`SELECT password_changed FROM %s WHERE username=$1 AND status='active';`, USERS_TABLE)

	var password_changed bool
	err := r.Db.QueryRow(context.Background(), query, username).Scan(&password_changed)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("user does not exist in table")
		return false, err
	}

	return password_changed, nil
}

// Retrieve public profile information by username
func (r *UsersRepository) GetUserProfile(username string) (*models.UserProfile, error) {
	query := fmt.Sprintf(`SELECT username, name, profile_photo, semester, karma FROM %s WHERE username=$1 AND status='active';`, USERS_TABLE)

	row, _ := r.Db.Query(context.Background(), query, username)
	profile, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserProfile])
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil, possible that user is not found")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &profile, nil
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
