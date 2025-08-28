package karma

// Represents the amount of points awarded for different interactions in the system
// There is no JSON binding for this. If the received setting is missing a field, 0 is set automatically.
type Karma struct {
	Semester         string `db:"semester" json:"semester"` // Does not matter what is in the request, configuration always updates for the latest semester only.
	CreateThreadPts  int    `db:"create_thread" json:"create_thread_pts"`
	CreateArticlePts int    `db:"create_article" json:"create_article_pts"`
	CreateCommentPts int    `db:"create_comment" json:"create_comment_pts"`
	CreatePostPts    int    `db:"create_post" json:"create_post_pts"`
	LikePts          int    `db:"receive_like" json:"like_pts"`
}
