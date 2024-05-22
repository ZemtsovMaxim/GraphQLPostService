package comments

type Comment struct {
	ID       int
	PostID   int
	Text     string
	ParentID *int // Может быть nil, если это корневой комментарий
}

type CommentRepository interface {
	CreateComment(comment *Comment) error
	GetCommentsByPostID(postID int, limit, offset int) ([]*Comment, error)
	GetReplies(parentID int) ([]*Comment, error)
	GetCommentByID(id int) (*Comment, error)
}
