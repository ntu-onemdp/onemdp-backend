package threads

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
)

type ThreadHandlers struct {
	NewThreadHandler    *NewThreadHandler
	DeleteThreadHandler *DeleteThreadHandler
}

func InitThreadHandlers(db *pgxpool.Pool) *ThreadHandlers {
	// Initialize repositories
	threadRepository := repositories.ThreadRepository{Db: db}
	postRepository := repositories.PostsRepository{Db: db}

	// Initialize services
	threadService := services.ThreadService{ThreadRepo: &threadRepository, PostRepo: &postRepository}

	// Initialize handlers
	newThreadHandler := NewThreadHandler{ThreadService: &threadService}
	deleteThreadHandler := DeleteThreadHandler{ThreadService: &threadService}

	return &ThreadHandlers{
		NewThreadHandler:    &newThreadHandler,
		DeleteThreadHandler: &deleteThreadHandler,
	}
}
