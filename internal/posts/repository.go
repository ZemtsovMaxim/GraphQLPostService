package posts

type Post struct {
	ID               int
	Title            string
	Content          string
	CommentsDisabled bool
}

type PostRepository interface {
	CreatePost(post *Post) error
	GetPostByID(id int) (*Post, error)
	GetAllPosts() ([]*Post, error)
	DisableComments(postID int) error
}
