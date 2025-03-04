package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type UserHandlers struct {
	UserProfileHandler    *ProfileHandler
	ChangePasswordHandler *ChangePasswordHandler
}

func InitUserHandlers(db *pgxpool.Pool) *UserHandlers {
	// Initialize repositories
	usersRepository := repositories.UsersRepository{Db: db}
	authRepository := repositories.AuthRepository{Db: db}

	// Initialize services
	userService := services.UserService{UsersRepo: &usersRepository}
	authService := services.AuthService{AuthRepo: &authRepository, UsersRepo: &usersRepository}

	// Initialize handlers
	profileHandler := ProfileHandler{UserService: &userService}
	changePasswordHandler := ChangePasswordHandler{AuthService: &authService}

	return &UserHandlers{
		UserProfileHandler:    &profileHandler,
		ChangePasswordHandler: &changePasswordHandler,
	}
}
