package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type UserHandlers struct {
	UserProfileHandler *ProfileHandler
}

func InitUserHandlers(db *pgxpool.Pool) *UserHandlers {
	// Initialize repositories
	usersRepository := repositories.UsersRepository{Db: db}

	// Initialize services
	userService := services.UserService{UsersRepo: &usersRepository}

	// Initialize handlers
	profileHandler := ProfileHandler{UserService: &userService}

	return &UserHandlers{
		UserProfileHandler: &profileHandler,
	}
}
