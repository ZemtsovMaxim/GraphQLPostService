package posts

import (
	"errors"
	"log/slog"
)

type InMemoryPostRepository struct {
	posts       map[int]*Post
	nextID      int
	commentsMap map[int]bool // Также отслеживаем посты с отключенными комментами
	logger      *slog.Logger
}

func NewInMemoryPostRepository(log *slog.Logger) *InMemoryPostRepository {
	return &InMemoryPostRepository{
		posts:       make(map[int]*Post),
		nextID:      1,
		commentsMap: make(map[int]bool),
		logger:      log,
	}
}

// Реализация методов интерфейса PostRepository для in-memory

func (r *InMemoryPostRepository) CreatePost(post *Post) error {
	post.ID = r.nextID
	r.posts[r.nextID] = post
	r.nextID++
	return nil
}

func (r *InMemoryPostRepository) GetPostByID(id int) (*Post, error) {
	post, exists := r.posts[id]
	if !exists {
		r.logger.Info("post not found")
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *InMemoryPostRepository) GetAllPosts() ([]*Post, error) {
	var result []*Post
	for _, post := range r.posts {
		result = append(result, post)
	}
	return result, nil
}

func (r *InMemoryPostRepository) DisableComments(postID int) error {
	if _, exists := r.posts[postID]; !exists {
		r.logger.Info("post not found")
		return errors.New("post not found")
	}
	r.posts[postID].CommentsDisabled = true
	return nil
}
