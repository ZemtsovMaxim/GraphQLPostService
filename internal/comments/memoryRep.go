package comments

import (
	"errors"
	"sync"
)

type InMemoryCommentRepository struct {
	mu       sync.RWMutex
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

	if limit == 0 && offset == 0 {
		return comments, nil
	}

	return comments[offset:end], nil
}

func (r *InMemoryCommentRepository) GetReplies(parentID int) ([]*Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var replies []*Comment
	for _, comment := range r.comments {
		for _, comm := range comment {
			if comm.ParentID != nil && *comm.ParentID == parentID {
				replies = append(replies, comm)
			}
		}
	}

	return replies, nil
}

func (r *InMemoryCommentRepository) GetCommentByID(id int) (*Comment, error) {
	for _, comment := range r.comments {
		for _, comm := range comment {
			if comm.ID == id {
				return comm, nil
			}
		}
	}
	return nil, nil
}
