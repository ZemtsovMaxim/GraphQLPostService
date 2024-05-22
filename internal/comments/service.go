package comments

import (
	"errors"
	"fmt"
)

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(postID int, text string, parentID *int) (*Comment, error) {
	comment := &Comment{PostID: postID, Text: text, ParentID: parentID}
	if *comment.ParentID != 0 {
		parentComment, err := s.repo.GetCommentByID(*comment.ParentID)
		if err != nil {
			return nil, fmt.Errorf("error fetching parent comment: %v", err)
		}
		if parentComment == nil {
			return nil, errors.New("parent comment does not exist")
		}
	}
	if len(text) > 1999 { // Ограничиваем длину комента
		panic("Comment len more than 2000 symbols")
	}
	err := s.repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID, limit, offset int) ([]*Comment, error) {
	return s.repo.GetCommentsByPostID(postID, limit, offset)
}
func (s *CommentService) GetLastCommentForPost(postID int) (*Comment, error) {
	// Получаем список комментариев к посту с помощью репозитория комментариев
	comments, err := s.GetCommentsByPostID(postID, 1, 0) // Получаем только один последний комментарий
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли комментарии к посту
	if len(comments) == 0 {
		return nil, fmt.Errorf("no comments found for post with ID %d", postID)
	}

	// Возвращаем первый (и единственный) комментарий из списка
	return comments[0], nil
}

func (s *CommentService) GetReplies(parentID int) ([]*Comment, error) {
	return s.repo.GetReplies(parentID)
}
