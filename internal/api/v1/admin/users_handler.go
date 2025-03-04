package admin

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type UsersHandler struct {
	CreateUserHandler      *CreateUserHandler
	GetUsersHandler        *GetUsersHandler
	UpdateUsersRoleHandler *UpdateUsersRoleHandler
}

func InitUserHandlers(db *pgxpool.Pool) *UsersHandler {
	// Initialize repositories
	usersRepository := repositories.UsersRepository{Db: db}
	authRepository := repositories.AuthRepository{Db: db}

	// Initialize services
	userService := services.UserService{UsersRepo: &usersRepository}
	authService := services.AuthService{AuthRepo: &authRepository, UsersRepo: &usersRepository}

	// Initialize handlers
	createUserHandler := CreateUserHandler{UserService: &userService, AuthService: &authService}
	getUsersHandler := GetUsersHandler{UserService: &userService}
	updateUsersRoleHandler := UpdateUsersRoleHandler{AuthService: &authService}

	return &UsersHandler{
		CreateUserHandler:      &createUserHandler,
		GetUsersHandler:        &getUsersHandler,
		UpdateUsersRoleHandler: &updateUsersRoleHandler,
	}
}
