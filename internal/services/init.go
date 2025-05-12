package services

import "github.com/ntu-onemdp/onemdp-backend/internal/repositories"

func Init() {
	Threads = NewThreadService(repositories.Threads, repositories.Posts, repositories.Likes)
	Posts = NewPostService(repositories.Posts)
	Likes = &LikeService{repositories.Likes}
	Auth = &AuthService{repositories.Auth, repositories.Users}
	Users = &UserService{repositories.Users}
	Images = &ImageService{repositories.Images}
}
