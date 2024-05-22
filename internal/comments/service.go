package comments

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(postID int, text string) (*Comment, error) {
	comment := &Comment{PostID: postID, Text: text}
	err := s.repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID, limit, offset int) ([]*Comment, error) {
	return s.repo.GetCommentsByPostID(postID, limit, offset)
}
