package comments

import (
	"errors"
)

type InMemoryCommentRepository struct {
	comments map[int][]*Comment
	nextID   int
}

func NewInMemoryCommentRepository() *InMemoryCommentRepository {
	return &InMemoryCommentRepository{
		comments: make(map[int][]*Comment),
		nextID:   1,
	}
}

// Реализация методов интерфейса CommentRepository для in-memory

func (r *InMemoryCommentRepository) CreateComment(comment *Comment) error {
	comment.ID = r.nextID
	r.comments[comment.PostID] = append(r.comments[comment.PostID], comment)
	r.nextID++
	return nil
}

func (r *InMemoryCommentRepository) GetCommentsByPostID(postID int, limit, offset int) ([]*Comment, error) {
	comments, exists := r.comments[postID]
	if !exists {
		return nil, errors.New("comments not found")
	}

	if offset > len(comments) {
		return []*Comment{}, nil
	}

	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}

	return comments[offset:end], nil
}
