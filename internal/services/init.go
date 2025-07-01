package services

import "github.com/ntu-onemdp/onemdp-backend/internal/repositories"

func Init() {
	Threads = NewThreadService(repositories.Threads, repositories.Posts, repositories.Likes)
	Posts = NewPostService(repositories.Posts)
	Likes = &LikeService{repositories.Likes}
	Users = &UserService{repositories.Users}
	Images = &ImageService{repositories.Images}
	Articles = NewArticleService(repositories.Articles, repositories.Comments)
	Comments = NewCommentService(repositories.Comments)
}
