package comments

type Comment struct {
	ID     int
	PostID int
	Text   string
}

type CommentRepository interface {
	CreateComment(comment *Comment) error
	GetCommentsByPostID(postID int, limit, offset int) ([]*Comment, error)
}
