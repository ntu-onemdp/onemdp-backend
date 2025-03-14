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
	LikePostHandlers  *LikePostHandlers
}

func InitPostHandlers(db *pgxpool.Pool) *PostHandlers {
	// Initialize repositories
	postRepository := repositories.PostsRepository{Db: db}
	threadRepository := repositories.ThreadRepository{Db: db}
	likeRepository := repositories.LikesRepository{Db: db}

	// Initialize services
	postService := services.NewPostService(&postRepository)
	threadService := services.NewThreadService(&threadRepository, &postRepository)
	likeService := services.NewLikeService(&likeRepository)

	// Initialize handlers
	newPostHandler := NewPostHandler{PostService: postService}
	updatePostHandler := UpdatePostHandler{PostService: postService, ThreadService: threadService}
	deletePostHandler := DeletePostHandler{PostService: postService}
	likePostHandlers := LikePostHandlers{likeService: likeService}

	return &PostHandlers{
		NewPostHandler:    &newPostHandler,
		DeletePostHandler: &deletePostHandler,
		UpdatePostHandler: &updatePostHandler,
		LikePostHandlers:  &likePostHandlers,
	}
}
