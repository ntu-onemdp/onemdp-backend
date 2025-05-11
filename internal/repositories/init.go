package repositories

import "github.com/jackc/pgx/v5/pgxpool"

// Initialize all repositories
func Init(db *pgxpool.Pool) {
	Auth = &AuthRepository{Db: db}
	Likes = &LikesRepository{Db: db}
	Posts = &PostsRepository{Db: db}
	Threads = &ThreadsRepository{Db: db}
	Users = &UsersRepository{Db: db}
}
