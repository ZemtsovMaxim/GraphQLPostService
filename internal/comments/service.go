package comments

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(postID int, text string) bool {
	comment := &Comment{
		PostID: postID,
		Text:   text,
	}
	if err := s.repo.CreateComment(comment); err != nil {
		return false
	}
	return true
}

func (s *CommentService) GetCommentsByPostID(postID int, limit, offset int) []*Comment {
	comments, err := s.repo.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		return nil
	}
	return comments
}
