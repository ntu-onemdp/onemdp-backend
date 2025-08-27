package karma

type Karma struct {
	Semester         string `db:"semester"`
	CreateThreadPts  int    `db:"create_thread" json:"create_thread_pts" required:"true"`
	CreateArticlePts int    `db:"create_article" json:"create_article_pts" required:"true"`
	CreateCommentPts int    `db:"create_comment" json:"create_comment_pts" required:"true"`
	CreatePostPts    int    `db:"create_post" json:"create_post_pts" required:"true"`
	LikePts          int    `db:"receive_like" json:"like_pts" required:"true"`
}
