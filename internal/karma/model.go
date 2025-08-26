package karma

type Karma struct {
	semester         string `db:"semester"`
	CreateThreadPts  int    `db:"create_thread"`
	CreateArticlePts int    `db:"create_article"`
	CreateCommentPts int    `db:"create_comment"`
	CreatePostPts    int    `db:"create_post"`
	LikePts          int    `db:"like"`
}
