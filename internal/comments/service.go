package comments

import (
	"errors"
	"log/slog"
)

type CommentService struct {
	repo   CommentRepository
	logger *slog.Logger
}

func NewCommentService(repo CommentRepository, log *slog.Logger) *CommentService {
	return &CommentService{repo: repo, logger: log}
}

func (s *CommentService) CreateComment(postID int, text string, parentID *int) (*Comment, error) {
	comment := &Comment{PostID: postID, Text: text, ParentID: parentID}
	if *comment.ParentID != 0 {
		parentComment, err := s.repo.GetCommentByID(*comment.ParentID)
		if err != nil {
			s.logger.Error("error fetching parent comment", slog.Any("err", err))
			return nil, err
		}
		if parentComment == nil {
			s.logger.Info("parent comment does not exist")
			return nil, errors.New("parent comment does not exist")
		}
	}
	if len(text) > 1999 { // Ограничиваем длину
		s.logger.Info("сomment len more than 2000 symbols")
		return nil, errors.New("сomment len more than 2000 symbols")
	}
	err := s.repo.CreateComment(comment)
	if err != nil {
		s.logger.Error("cant create comment", slog.Any("err", err))
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID, limit, offset int) ([]*Comment, error) {
	return s.repo.GetCommentsByPostID(postID, limit, offset)
}
func (s *CommentService) GetLastCommentForPost(postID int) (*Comment, error) {
	comments, err := s.GetCommentsByPostID(postID, 1, 0) // Получаем только один последний комментарий
	if err != nil {
		s.logger.Error("error get comment by post", slog.Any("err", err))
		return nil, err
	}

	// Проверяем, есть ли комментарии к посту
	if len(comments) == 0 {
		s.logger.Info("no comments found for post with ID")
		return nil, errors.New("no comments found for post with ID")
	}

	// Возвращаем первый (и единственный) комментарий из списка
	return comments[0], nil
}

func (s *CommentService) GetReplies(parentID int) ([]*Comment, error) {
	return s.repo.GetReplies(parentID)
}
