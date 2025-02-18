package main

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/admin"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
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
	userHandler := admin.CreateUserHandler{UserService: &userService}

	// Register routes
	routes.RegisterLoginRoute(r, &authHandler)
	routes.RegisterCreateUsersRoute(r, &userHandler)

	// // Protect admin routes
	// protected := r.Group("/api/v1/admin", auth.AdminGuard())

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 		"content": "hello",
	// 	})
	// })

	// r.POST("/api/v1/auth/login", func(c *gin.Context) {
	// 	auth.HandleLogin(c, db.Pool)
	// })

	// // Admin functions
	// // Verify if user is admin. Used in frontend to grant access to protected routes.
	// protected.GET("/verify", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"role": "admin",
	// 	})
	// })

	// // Enrol new users
	// protected.POST("/users/create", func(c *gin.Context) {
	// 	users.CreateUsers(c, db.Pool)
	// })
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
