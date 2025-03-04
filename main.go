package main

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/api/middlewares"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/users"
	"github.com/ntu-onemdp/onemdp-backend/internal/db"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"

	"github.com/gin-gonic/gin"
	routes "github.com/ntu-onemdp/onemdp-backend/internal/api"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	db.Init()
	defer db.Close()

	r := gin.Default()

	r.Use(cors.Default())

	// Initialize repositories
	authRepo := repositories.AuthRepository{Db: db.Pool}
	usersRepo := repositories.UsersRepository{Db: db.Pool}

	// Initialize services
	authService := services.AuthService{AuthRepo: &authRepo, UsersRepo: &usersRepo}
	userService := services.UserService{UsersRepo: &usersRepo}

	// Initialize handlers (might be shifted in the future)
	authHandler := auth.LoginHandler{AuthService: &authService}
	profileHandler := users.ProfileHandler{UserService: &userService}

	// Register public routes
	routes.RegisterLoginRoute(r, &authHandler)

	// Register student routes
	studentRoutes := r.Group("/api/v1/users/:username", middlewares.AuthGuard())
	routes.RegisterStudentUserRoutes(studentRoutes, &profileHandler)

	// Register admin routes
	adminRoutes := r.Group("/api/v1/admin", middlewares.AdminGuard())
	routes.RegisterAdminUserRoutes(adminRoutes, db.Pool)

	// Ping route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"content": "hello",
		})
	})

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
