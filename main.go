package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ntu-onemdp/onemdp-backend/internal/api/middlewares"
	"github.com/ntu-onemdp/onemdp-backend/internal/api/v1/auth"
	"github.com/ntu-onemdp/onemdp-backend/internal/db"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"

	"github.com/gin-gonic/gin"
	routes "github.com/ntu-onemdp/onemdp-backend/internal/api"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	db.Init()
	defer db.Close()

	// Initialize JWT handler
	utils.InitJwt()

	r := gin.Default()

	r.Use(cors.Default())

	// Initialize repositories
	authRepo := repositories.AuthRepository{Db: db.Pool}
	usersRepo := repositories.UsersRepository{Db: db.Pool}

	// Initialize services
	authService := services.AuthService{AuthRepo: &authRepo, UsersRepo: &usersRepo}

	// Initialize handlers (might be shifted in the future)
	authHandler := auth.LoginHandler{AuthService: &authService}

	// Register public routes
	routes.RegisterLoginRoute(r, &authHandler)

	// Register change password route
	routes.RegisterChangePasswordRoute(r, db.Pool)

	// Register student routes
	studentRoutes := r.Group("/api/v1/users/:username", middlewares.AuthGuard())
	routes.RegisterStudentUserRoutes(studentRoutes, db.Pool)

	// Register thread routes
	threadRoutes := r.Group("/api/v1/threads", middlewares.AuthGuard())
	routes.RegisterThreadRoutes(threadRoutes, db.Pool)

	// Register post routes
	postRoutes := r.Group("/api/v1/posts", middlewares.AuthGuard())
	routes.RegisterPostRoutes(postRoutes, db.Pool)

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

	// For debugging purposes only
	r.POST("/ping", func(c *gin.Context) {
		var requestBody map[string]interface{} // Use a generic map to handle any JSON structure

		// Bind the JSON body to the requestBody variable
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		fmt.Println("Request Body:", requestBody)

		timestamp := time.Now().Unix()
		// Echo back the received JSON
		c.JSON(http.StatusOK, gin.H{
			"body":      requestBody,
			"timestamp": timestamp,
		})
	})

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
