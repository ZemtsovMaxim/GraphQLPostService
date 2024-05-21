package posts

import (
	"database/sql"
)

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) *PostgresPostRepository {
	return &PostgresPostRepository{db: db}
}

// Реализация методов интерфейса PostRepository для PostgreSQL

func (r *PostgresPostRepository) CreatePost(post *Post) error {
	_, err := r.db.Exec("INSERT INTO posts (title, content, comments_disabled) VALUES ($1, $2, $3)", post.Title, post.Content, post.CommentsDisabled)
	return err
}

func (r *PostgresPostRepository) GetPostByID(id int) (*Post, error) {
	post := &Post{}
	err := r.db.QueryRow("SELECT id, title, content, comments_disabled FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostgresPostRepository) GetAllPosts() ([]*Post, error) {
	rows, err := r.db.Query("SELECT id, title, content, comments_disabled FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostgresPostRepository) DisableComments(postID int) error {
	_, err := r.db.Exec("UPDATE posts SET comments_disabled = TRUE WHERE id = $1", postID)
	return err
}
