package posts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryPostRepository_CreatePost(t *testing.T) {
	// Создаем новый репозиторий
	repo := NewInMemoryPostRepository(nil)

	// Создаем пост
	post := &Post{Title: "Test Post", Content: "Test Content"}
	err := repo.CreatePost(post)

	// Проверяем, что метод создания поста не вернул ошибку
	require.NoError(t, err)

	// Проверяем, что пост был добавлен в репозиторий с правильным ID
	assert.Equal(t, 1, post.ID)
}

func TestInMemoryPostRepository_GetPostByID(t *testing.T) {
	// Создаем новый репозиторий
	repo := NewInMemoryPostRepository(nil)

	// Создаем и добавляем пост в репозиторий
	post := &Post{ID: 1, Title: "Test Post", Content: "Test Content"}
	repo.posts[1] = post

	// Получаем пост по ID
	foundPost, err := repo.GetPostByID(1)

	// Проверяем, что метод получения поста не вернул ошибку
	require.NoError(t, err)

	// Проверяем, что полученный пост совпадает с ожидаемым
	assert.Equal(t, post, foundPost)
}

func TestInMemoryPostRepository_GetAllPosts(t *testing.T) {
	// Создаем новый репозиторий
	repo := NewInMemoryPostRepository(nil)

	// Создаем и добавляем посты в репозиторий
	posts := []*Post{
		{ID: 1, Title: "Test Post 1", Content: "Test Content 1"},
		{ID: 2, Title: "Test Post 2", Content: "Test Content 2"},
	}
	for _, post := range posts {
		repo.posts[post.ID] = post
	}

	// Получаем все посты из репозитория
	allPosts, err := repo.GetAllPosts()

	// Проверяем, что метод получения всех постов не вернул ошибку
	require.NoError(t, err)

	// Проверяем, что количество полученных постов совпадает с ожидаемым
	assert.Len(t, allPosts, len(posts))

	// Проверяем, что каждый полученный пост совпадает с ожидаемым
	for _, expectedPost := range posts {
		assert.Contains(t, allPosts, expectedPost)
	}
}

func TestInMemoryPostRepository_DisableComments(t *testing.T) {
	// Создаем новый репозиторий
	repo := NewInMemoryPostRepository(nil)

	// Создаем и добавляем пост в репозиторий
	post := &Post{ID: 1, Title: "Test Post", Content: "Test Content"}
	repo.posts[1] = post

	// Отключаем комментарии для поста
	err := repo.DisableComments(1)

	// Проверяем, что метод отключения комментариев не вернул ошибку
	require.NoError(t, err)

	// Проверяем, что комментарии действительно отключены для поста
	assert.True(t, post.CommentsDisabled)
}
