package posts

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type PostHandlers struct {
	NewPostHandler    *NewPostHandler
	DeletePostHandler *DeletePostHandler
	UpdatePostHandler *UpdatePostHandler
}

func InitPostHandlers(db *pgxpool.Pool) *PostHandlers {
	// Initialize repositories
	postRepository := repositories.PostsRepository{Db: db}
	threadRepository := repositories.ThreadRepository{Db: db}

	// Initialize services
	postService := services.PostService{PostRepo: &postRepository}
	threadService := services.ThreadService{ThreadRepo: &threadRepository, PostRepo: &postRepository}

	// Initialize handlers
	newPostHandler := NewPostHandler{PostService: &postService}
	updatePostHandler := UpdatePostHandler{PostService: &postService, ThreadService: &threadService}
	deletePostHandler := DeletePostHandler{PostService: &postService}

	return &PostHandlers{
		NewPostHandler:    &newPostHandler,
		DeletePostHandler: &deletePostHandler,
		UpdatePostHandler: &updatePostHandler,
	}
}
