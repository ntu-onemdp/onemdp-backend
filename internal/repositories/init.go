package repositories

import "github.com/jackc/pgx/v5/pgxpool"

// Initialize all repositories
func Init(db *pgxpool.Pool) {
	Likes = &LikesRepository{Db: db}
	Posts = &PostsRepository{Db: db}
	Threads = &ThreadsRepository{Db: db}
	Users = &UsersRepository{Db: db}
	Images = &ImagesRepository{Db: db}
	Articles = &ArticleRepository{Db: db}
	Comments = &CommentsRepository{Db: db}
}
