package posts

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

// GetAllPosts возвращает список всех постов.
func (s *PostService) GetAllPosts() ([]*Post, error) {
	return s.repo.GetAllPosts()
}

// GetPostByID возвращает пост с указанным идентификатором.
func (s *PostService) GetPostByID(id int) (*Post, error) {
	return s.repo.GetPostByID(id)
}

// CreatePost создает новый пост с указанным заголовком и содержимым.
func (s *PostService) CreatePost(title, content string) error {
	post := &Post{
		Title:   title,
		Content: content,
	}
	return s.repo.CreatePost(post)
}

// DisableComments запрещает оставление комментариев к посту с указанным идентификатором.
func (s *PostService) DisableComments(postID int) error {
	return s.repo.DisableComments(postID)
}
