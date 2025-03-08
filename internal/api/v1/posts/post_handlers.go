package posts

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type PostHandlers struct {
	NewPostHandler *NewPostHandler
}

func InitPostHandlers(db *pgxpool.Pool) *PostHandlers {
	// Initialize repositories
	postRepository := repositories.PostsRepository{Db: db}

	// Initialize services
	postService := services.PostService{PostRepo: &postRepository}

	// Initialize handlers
	newPostHandler := NewPostHandler{PostService: &postService}

	return &PostHandlers{
		NewPostHandler: &newPostHandler,
	}
}
