package posts

import "log/slog"

type PostService struct {
	repo   PostRepository
	logger *slog.Logger
}

func NewPostService(repo PostRepository, log *slog.Logger) *PostService {
	return &PostService{repo: repo, logger: log}
}

func (s *PostService) CreatePost(title, content string) (*Post, error) {
	post := &Post{Title: title, Content: content}
	err := s.repo.CreatePost(post)
	if err != nil {
		s.logger.Error("cant create post", slog.Any("err", err))
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
		s.logger.Error("cant disable comments ", slog.Any("err", err))
		return false, err
	}
	return true, nil
}
