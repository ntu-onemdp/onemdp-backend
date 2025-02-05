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
	query := "SELECT status FROM users WHERE username=$1"

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
	query := "SELECT * FROM users WHERE username=$1 AND status = 'active'"

	row, _ := r.Db.Query(context.Background(), query, username)
	user, err := pgx.RowToStructByName[models.User](row)

	if err != nil {
		utils.Logger.Error().Err(err)
		return nil, err
	}

	return &user, nil
}
