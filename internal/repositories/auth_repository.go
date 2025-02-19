package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type AuthRepository struct {
	Db *pgxpool.Pool
}

func (r *AuthRepository) GetAuthByUsername(username string) (*models.AuthModel, error) {
	query := "SELECT * FROM auth WHERE username=$1"

	var auth models.AuthModel
	err := r.Db.QueryRow(context.Background(), query, username).Scan(&auth.Username, &auth.Password, &auth.Role)
	if err != nil {
		utils.Logger.Error().Err(err)
		return nil, err
	}

	return &auth, nil
}
