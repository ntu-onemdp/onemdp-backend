package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// Users table name in db
const USERS_TABLE = "users"
const PENDING_USERS_TABLE = "pending_users"

type UsersRepository struct {
	Db *pgxpool.Pool
}

var Users *UsersRepository

// InsertOneUser inserts a new pending user into the database after verifying email uniqueness.
//
// This repository method performs an idempotent insert into the pending_users table,
// ensuring the email doesn't already exist in the main users table. Designed for use
// during user registration workflows.
//
// Parameters:
//   - user: Pointer to PendingUser model containing email, role, and semester
//
// Returns:
//   - error: nil on successful insertion, error in these cases:
//   - Database connection/query errors
//   - Constraint violations (e.g. invalid email format)
//   - Email already exists in users table
//
// Example usage:
//
//	newUser := &models.PendingUser{
//	    Email: "test@example.com",
//	    Role: "student",
//	    Semester: 2,
//	}
//	err := usersRepo.InsertOneUser(newUser)
func (r *UsersRepository) InsertOneUser(user *models.PendingUser) error {
	query := `
	INSERT INTO pending_users (email, role, semester) 
	VALUES ($1, $2, $3)
	WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = $1);` // Ensure that the user does not already exist in the users table

	_, err := r.Db.Exec(context.Background(), query, user.Email, user.Role, user.Semester)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg(fmt.Sprintf("%s successfully inserted into database", user.Email))
	return nil
}

// RegisterUser finalizes a user's registration by moving them from the pending_users table to the users table.
//
// Parameters:
//   - uid:   The unique identifier for the user (e.g., from authentication provider).
//   - email: The email address of the user to register (must exist in pending_users).
//   - name:  The full name of the user.
//
// Returns:
//   - error: Returns an error if any step fails (e.g., transaction errors, user not found in pending_users, database issues).
func (r *UsersRepository) RegisterUser(uid string, email string, name string) error {
	ctx := context.Background()

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)

	// Retrieve pending user
	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1;`, PENDING_USERS_TABLE)

	row, _ := tx.Query(ctx, query, email)
	defer row.Close()
	pending_user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.PendingUser])
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to retrieve pending user")
		return err
	}
	utils.Logger.Trace().Msgf("Retrieved pending user: %s", pending_user.Email)

	// Create new user from pending user
	user := models.CreateUser(uid, name, pending_user.Email, pending_user.Semester, pending_user.Role)

	// Insert user into users table
	query = fmt.Sprintf(`
	INSERT INTO %s (email, uid, name, role, semester)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (email) DO UPDATE
	SET uid = EXCLUDED.uid, name = EXCLUDED.name, role = EXCLUDED.role, semester = EXCLUDED.semester;`, USERS_TABLE)

	if _, err := tx.Exec(ctx, query, user.Email, user.Uid, user.Name, user.Role, user.Semester); err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to insert user into users table")
		return err
	}
	utils.Logger.Trace().Msgf("Inserted user: %s into users table", user.Email)

	// Delete pending user
	query = fmt.Sprintf(`DELETE FROM %s WHERE email=$1;`, PENDING_USERS_TABLE)
	if _, err := tx.Exec(ctx, query, pending_user.Email); err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to delete pending user")
		return err
	}
	utils.Logger.Trace().Msgf("Deleted pending user: %s", pending_user.Email)

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to commit transaction")
	}

	utils.Logger.Debug().Msgf("User %s successfully registered", user.Email)
	return nil
}

// Dangerously insert user directly into database
// Used to insert Eduvisor bot in
func (r *UsersRepository) DangerouslyInsertUser(user *models.User) error {
	query := fmt.Sprintf(`INSERT INTO %s (UID, NAME, EMAIL, ROLE, SEMESTER) VALUES ($1, $2, $3, $4, $5);`, USERS_TABLE)

	if _, err := r.Db.Exec(context.Background(), query, user.Uid, user.Name, user.Email, user.Role, user.Semester); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error inserting user with uid %s", user.Uid)
		return err
	}

	utils.Logger.Info().Msgf("User with uid %s and name %s successfully inserted into %s", user.Uid, user.Name, USERS_TABLE)
	return nil
}

// Checks if user is pending registration
func (r *UsersRepository) IsUserPending(email string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email=$1);`, PENDING_USERS_TABLE)
	var exists bool
	err := r.Db.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking if user is pending")
		return false, err
	}

	return exists, nil
}

// Check if Eduvisor is registered. If Eduvisor is registered, return user object. Else, return nil.
func (r *UsersRepository) GetEduvisor() *models.User {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE NAME=$1;`, USERS_TABLE)

	row, _ := r.Db.Query(context.Background(), query, constants.EDUVISOR_NAME)
	eduvisor, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("")
		utils.Logger.Warn().Msg("Eduvisor not found in the system, returning nil")
		return nil
	}

	utils.Logger.Debug().Msg("Eduvisor found in user system")
	return eduvisor
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

// Retrieve user's status using uid.
// Throws error if uid cannot be found
// Deprecate in the future
func (r *UsersRepository) GetStatus(uid string) (string, error) {
	query := fmt.Sprintf(`SELECT status FROM %s WHERE uid=$1;`, USERS_TABLE)

	var status string
	err := r.Db.QueryRow(context.Background(), query, uid).Scan(&status)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return "", err
	}

	return status, nil
}

// Retrieve user's information from uid
// This function *checks* if user is active before returning. If the user's status is not 'active',
// an error is return instead.
func (r *UsersRepository) GetUserByUid(uid string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE uid=$1 AND status='active';", USERS_TABLE)
	row, _ := r.Db.Query(context.Background(), query, uid)
	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	utils.Logger.Trace().Interface("user", user).Msg("")
	return user, nil
}

// Admin: Retrieve user's information from uid
// This function is able to retrieve deleted users as well
func (r *UsersRepository) GetUserByUidAdmin(uid string) (*models.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE uid=$1;`, USERS_TABLE)
	row, _ := r.Db.Query(context.Background(), query, uid)
	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	return user, nil
}

// Retrieve public profile information by uid
func (r *UsersRepository) GetUserProfile(uid string) (*models.UserProfile, error) {
	query := fmt.Sprintf(`SELECT uid, email, name, role, profile_photo, semester, karma FROM %s WHERE uid=$1 AND status='active';`, USERS_TABLE)

	row, _ := r.Db.Query(context.Background(), query, uid)
	profile, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserProfile])
	if err != nil {
		utils.Logger.Debug().Msg("Returning nil, possible that user is not found")
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	return &profile, nil
}

// Retrieve user's profile photo from database
// Do not filter for active users only
func (r *UsersRepository) GetProfilePhoto(uid string) ([]byte, error) {
	query := fmt.Sprintf(`SELECT PROFILE_PHOTO FROM %s WHERE UID=$1;`, USERS_TABLE)

	var image []byte
	if err := r.Db.QueryRow(context.Background(), query, uid).Scan(&image); err != nil {
		utils.Logger.Warn().Err(err).Msgf("User of uid %s not found.", uid)
		return nil, err
	}

	utils.Logger.Debug().Msgf("Retrieved profile photo for uid %s", uid)
	return image, nil
}

// Retrieve user's role
func (r *UsersRepository) GetUserRole(uid string) (string, error) {
	query := fmt.Sprintf(`SELECT role FROM %s WHERE uid=$1;`, USERS_TABLE)

	var role string
	err := r.Db.QueryRow(context.Background(), query, uid).Scan(&role)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error retrieving user role")
		return "", err
	}

	utils.Logger.Debug().Msgf("User %s has role %s", uid, role)
	return role, nil
}

// Retrieve rankings of top N users by karma and current semester.
func (r *UsersRepository) GetTopKarma(semester string, n int) ([]models.UserProfile, error) {
	query := fmt.Sprintf(`SELECT uid, email, name, role, profile_photo, semester, karma FROM %s WHERE ROLE='student' AND STATUS='active' AND SEMESTER=$1 ORDER BY KARMA DESC LIMIT $2;`, USERS_TABLE)

	rows, _ := r.Db.Query(context.Background(), query, semester, n)
	profiles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserProfile])
	if err != nil {
		utils.Logger.Error().Err(err).Msgf("Error retrieving profiles for semester %s", semester)
		return nil, err
	}

	if len(profiles) == 0 {
		utils.Logger.Warn().Msgf("0 profiles retrieved for semester %s", semester)
	} else {
		utils.Logger.Debug().Msgf("%d profiles retrieved for semester %s", len(profiles), semester)
	}

	return profiles, nil
}

// Update user's role
func (r *UsersRepository) UpdateUserRole(uid string, role string) error {
	query := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE uid=$2;`, USERS_TABLE)

	_, err := r.Db.Exec(context.Background(), query, role, uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating user role")
		return err
	}

	utils.Logger.Trace().Msgf("User %s role updated to %s", uid, role)
	return nil
}

// Update user profile photo. Returns nil on success.
func (r *UsersRepository) UpdateProfilePhoto(uid string, image []byte) error {
	query := fmt.Sprintf(`UPDATE %s SET PROFILE_PHOTO=$1 WHERE UID=$2;`, USERS_TABLE)

	if _, err := r.Db.Exec(context.Background(), query, image, uid); err != nil {
		utils.Logger.Error().Err(err).Msgf("Error updating profile photo for user %s", uid)
		return err
	}

	utils.Logger.Info().Msgf("Successfully updated profile photo for user %s", uid)
	return nil
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
