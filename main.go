package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ntu-onemdp/onemdp-backend/internal/api/middlewares"
	"github.com/ntu-onemdp/onemdp-backend/internal/db"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
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
	services.InitJwt()

	r := gin.Default()

	// Reduce max memory limit for multipart form data
	r.MaxMultipartMemory = 4 << 20 // 4 MiB

	r.Use(cors.Default())

	// Initialize repositories
	repositories.Init(db.Pool)

	// Initialize services
	services.Init()

	// Initialize eduvisor service
	services.Eduvisor = services.NewEduvisorService()

	// Register public routes
	routes.RegisterLoginRoute(r)

	// Register student routes
	studentRoutes := r.Group("/api/v1/users", middlewares.AuthGuard(models.Student))
	routes.RegisterStudentUserRoutes(studentRoutes)

	// Register thread routes
	threadRoutes := r.Group("/api/v1/threads", middlewares.AuthGuard(models.Student))
	routes.RegisterThreadRoutes(threadRoutes)

	// Register post routes
	postRoutes := r.Group("/api/v1/posts", middlewares.AuthGuard(models.Student))
	routes.RegisterPostRoutes(postRoutes)

	// Register article routes
	articleRoutes := r.Group("/api/v1/articles", middlewares.AuthGuard(models.Student))
	routes.RegisterArticleRoutes(articleRoutes)

	// Register comment routes
	commentRoutes := r.Group("/api/v1/comments", middlewares.AuthGuard(models.Student))
	routes.RegisterCommentRoutes(commentRoutes)

	// Register image routes
	imageRoutes := r.Group("/api/v1/images", middlewares.AuthGuard(models.Student))
	routes.RegisterImageRoutes(imageRoutes)

	// Register like content routes
	likeRoutes := r.Group("/api/v1/like", middlewares.AuthGuard(models.Student))
	routes.RegisterLikeRoutes(likeRoutes)

	// Register favorite content routes
	favoriteRoutes := r.Group("/api/v1/saved", middlewares.AuthGuard(models.Student))
	routes.RegisterSavedRoutes(favoriteRoutes)

	// Register student file routes
	fileRoutes := r.Group("/api/v1/files", middlewares.AuthGuard(models.Student))
	routes.RegisterFileRoutes(fileRoutes)

	// Register staff file routes
	staffFileRoutes := r.Group("/api/v1/staff/files", middlewares.AuthGuard(models.Staff))
	routes.RegisterFileMgmtRoutes(staffFileRoutes)

	// Register admin routes
	adminRoutes := r.Group("/api/v1/admin", middlewares.AuthGuard(models.Admin))
	routes.RegisterAdminUserRoutes(adminRoutes)

	utils.Logger.Warn().Msg("/ping routes are active. Remove them for production")

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
