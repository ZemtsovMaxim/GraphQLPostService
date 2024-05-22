package posts

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(title, content string) (*Post, error) {
	post := &Post{Title: title, Content: content}
	err := s.repo.CreatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByID(id int) (*Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *PostService) GetAllPosts() ([]*Post, error) {
	return s.repo.GetAllPosts()
}

func (s *PostService) DisableComments(id int) (bool, error) {
	err := s.repo.DisableComments(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
