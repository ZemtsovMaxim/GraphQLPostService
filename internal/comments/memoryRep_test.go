package comments

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommentRepository_CreateComment(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	repo := NewInMemoryCommentRepository(logger)

	comment := &Comment{
		Text:   "Test Comment",
		PostID: 1,
	}

	err := repo.CreateComment(comment)
	assert.NoError(t, err)

	// Убеждаемся, что идентификатор комментария присвоен
	assert.NotZero(t, comment.ID)

}

func TestInMemoryCommentRepository_GetCommentsByPostID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	repo := NewInMemoryCommentRepository(logger)

	comment1 := &Comment{
		Text:   "Test Comment 1",
		PostID: 1,
	}
	comment2 := &Comment{
		Text:   "Test Comment 2",
		PostID: 2,
	}

	err := repo.CreateComment(comment1)

	if err != nil {
		logger.Info("error create comment1")
	}

	err2 := repo.CreateComment(comment2)

	if err2 != nil {
		logger.Info("error create comment2")
	}

	// Тестируем, когда комментарии существуют для заданного идентификатора сообщения
	comments, err := repo.GetCommentsByPostID(1, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Test Comment 1", comments[0].Text)

	// Тестируем, когда комментарии не существуют для заданного идентификатора сообщения
	_, err = repo.GetCommentsByPostID(3, 0, 0)
	assert.Error(t, err)
}

func TestInMemoryCommentRepository_GetReplies(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	repo := NewInMemoryCommentRepository(logger)

	comment1 := &Comment{
		ID:       1,
		Text:     "Test Comment 1",
		PostID:   1,
		ParentID: nil,
	}
	comment2 := &Comment{
		ID:       2,
		Text:     "Reply to Comment 1",
		PostID:   1,
		ParentID: &comment1.ID,
	}
	comment3 := &Comment{
		ID:       3,
		Text:     "Test Comment 2",
		PostID:   2,
		ParentID: nil,
	}

	err := repo.CreateComment(comment1)

	if err != nil {
		logger.Info("error create comment1")
	}

	err2 := repo.CreateComment(comment2)

	if err2 != nil {
		logger.Info("error create comment2")
	}
	err3 := repo.CreateComment(comment3)

	if err3 != nil {
		logger.Info("error create comment3")
	}

	// Тестируем получение ответов на комментарий с ID = 1
	replies, err := repo.GetReplies(comment1.ID)
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "Reply to Comment 1", replies[0].Text)

	// Тестируем, когда ответов нет
	replies, err = repo.GetReplies(comment3.ID)
	assert.NoError(t, err)
	assert.Len(t, replies, 0)

}

func TestInMemoryCommentRepository_GetCommentByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	repo := NewInMemoryCommentRepository(logger)

	comment := &Comment{
		ID:     1,
		Text:   "Test Comment",
		PostID: 1,
	}

	err := repo.CreateComment(comment)
	if err != nil {
		logger.Info("error create comment")
	}

	// Тестируем получение комментария по его идентификатору
	result, err := repo.GetCommentByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Comment", result.Text)

	// Тестируем, когда комментарий не найден
	result, err = repo.GetCommentByID(2)
	assert.NoError(t, err)
	assert.Nil(t, result)

}
