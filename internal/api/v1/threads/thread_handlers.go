package threads

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type ThreadHandlers struct {
	NewThreadHandler    *CreateThreadHandler
	DeleteThreadHandler *DeleteThreadHandler
	GetThreadHandler    *GetThreadHandler
	LikeThreadHandlers  *LikeThreadHandlers
}

func InitThreadHandlers(db *pgxpool.Pool) *ThreadHandlers {
	// Initialize repositories
	threadRepository := repositories.ThreadRepository{Db: db}
	postsRepository := repositories.PostsRepository{Db: db}
	likesRepository := repositories.LikesRepository{Db: db}

	// Initialize services
	threadService := services.NewThreadService(&threadRepository, &postsRepository)
	likeService := services.NewLikeService(&likesRepository)

	// Initialize handlers
	newThreadHandler := CreateThreadHandler{ThreadService: threadService}
	deleteThreadHandler := DeleteThreadHandler{ThreadService: threadService}
	getThreadHandler := GetThreadHandler{ThreadService: threadService}
	likeThreadHandlers := LikeThreadHandlers{likeService: likeService, threadService: threadService}

	return &ThreadHandlers{
		NewThreadHandler:    &newThreadHandler,
		DeleteThreadHandler: &deleteThreadHandler,
		GetThreadHandler:    &getThreadHandler,
		LikeThreadHandlers:  &likeThreadHandlers,
	}
}
