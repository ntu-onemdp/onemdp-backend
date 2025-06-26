package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ntu-onemdp/onemdp-backend/internal/api/middlewares"
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

	// Initialize JWT handler
	services.InitJwt()

	r := gin.Default()

	// Reduce max memory limit for multipart form data
	r.MaxMultipartMemory = 4 << 20 // 4 MiB

	r.Use(cors.Default())

	// Initialize repositories
	repositories.Init(db.Pool)

	// Initialize services
	services.Init()

	// Register public routes
	routes.RegisterLoginRoute(r)

	// Register student routes
	studentRoutes := r.Group("/api/v1/users", middlewares.AuthGuard())
	routes.RegisterStudentUserRoutes(studentRoutes)

	// Register thread routes
	threadRoutes := r.Group("/api/v1/threads", middlewares.AuthGuard())
	routes.RegisterThreadRoutes(threadRoutes)

	// Register post routes
	postRoutes := r.Group("/api/v1/posts", middlewares.AuthGuard())
	routes.RegisterPostRoutes(postRoutes)

	// Register article routes
	articleRoutes := r.Group("/api/v1/articles", middlewares.AuthGuard())
	routes.RegisterArticleRoutes(articleRoutes)

	// Register image routes
	imageRoutes := r.Group("/api/v1/images", middlewares.AuthGuard())
	routes.RegisterImageRoutes(imageRoutes)

	// Register admin routes
	adminRoutes := r.Group("/api/v1/admin", middlewares.AdminGuard())
	routes.RegisterAdminUserRoutes(adminRoutes)

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
