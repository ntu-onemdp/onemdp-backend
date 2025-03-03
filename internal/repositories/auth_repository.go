package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type AuthRepository struct {
	Db *pgxpool.Pool
}

const AUTH_TABLE = "auth"

// Insert new auth detail
func (r *AuthRepository) InsertAuthDetails(auth *models.AuthModel) error {
	query := fmt.Sprintf(`INSERT INTO %s (username, password, role) VALUES ($1, $2, $3);`, AUTH_TABLE)

	_, err := r.Db.Exec(context.Background(), query, auth.Username, auth.Password, auth.Role)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg("Successfully inserted auth for " + auth.Username)
	return nil
}

// Retrieve user auth details using usermame
func (r *AuthRepository) GetAuthByUsername(username string) (*models.AuthModel, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1;", AUTH_TABLE)

	var auth models.AuthModel
	err := r.Db.QueryRow(context.Background(), query, username).Scan(&auth.Username, &auth.Password, &auth.Role)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return nil, err
	}

	utils.Logger.Trace().Msg("Successfully retrieved auth for " + username)
	return &auth, nil
}

// Update individual user role. Returns nil on success
func (r *AuthRepository) UpdateUserRole(username string, new_role string) error {
	query := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE username=$2;`, AUTH_TABLE)

	_, err := r.Db.Exec(context.Background(), query, new_role, username)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("")
		return err
	}

	utils.Logger.Trace().Msg("Successfully updated role for " + username)
	return nil
}
