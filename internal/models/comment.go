package models

// Comment models how a comment is represented on the API.
type Comment struct {
	DbComment

	Author   string `json:"author" db:"author_name"` // Name of the author
	NumLikes int    `json:"num_likes" db:"num_likes"`
	IsLiked  bool   `json:"is_liked" db:"is_liked"` // Whether the post has been liked by user
}

// DbComment models how a comment is stored in the database
type DbComment struct {
	CommentID   string `json:"comment_id" db:"comment_id"`
	AuthorUID   string `json:"author_uid" db:"author" binding:"required"`
	ArticleID   string `json:"article_id" db:"article_id" `
	Content     string `json:"content" db:"content" binding:"required"`
	TimeCreated string `json:"time_created" db:"time_created"`
	LastEdited  string `json:"last_edited" db:"last_edited"`
	Flagged     bool   `json:"flagged" db:"flagged"`
	IsAvailable bool   `json:"is_available" db:"is_available"`
}

type CommentFactory struct {
	ContentFactory
}

func NewCommentFactory() *CommentFactory {
	return &CommentFactory{}
}

func (f *CommentFactory) New(authorUID string, articleID string) {

}
