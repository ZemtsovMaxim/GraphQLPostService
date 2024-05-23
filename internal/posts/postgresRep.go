package posts

import (
	"database/sql"
	"log/slog"
)

type PostgresPostRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresPostRepository(db *sql.DB, log *slog.Logger) *PostgresPostRepository {
	return &PostgresPostRepository{db: db, logger: log}
}

// Реализация методов интерфейса PostRepository для PostgreSQL

func (r *PostgresPostRepository) CreatePost(post *Post) error {
	_, err := r.db.Exec("INSERT INTO posts (title, content, comments_disabled) VALUES ($1, $2, $3)", post.Title, post.Content, post.CommentsDisabled)
	if err != nil {
		r.logger.Error("Cant create post", slog.Any("err", err))
	}
	return err
}

func (r *PostgresPostRepository) GetPostByID(id int) (*Post, error) {
	row := r.db.QueryRow("SELECT id, title, content, comments_disabled FROM posts WHERE id = $1", id)
	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
	if err != nil {
		r.logger.Error("Cant get post by id", slog.Any("err", err))
		return nil, err
	}
	return &post, nil
}

func (r *PostgresPostRepository) GetAllPosts() ([]*Post, error) {
	rows, err := r.db.Query("SELECT id, title, content, comments_disabled FROM posts")
	if err != nil {
		r.logger.Error("Cant get all posts", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled); err != nil {
			r.logger.Error("Cant scan all posts", slog.Any("err", err))
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostgresPostRepository) DisableComments(postID int) error {
	_, err := r.db.Exec("UPDATE posts SET comments_disabled = TRUE WHERE id = $1", postID)
	if err != nil {
		r.logger.Error("Cant disable comments", slog.Any("err", err))
	}
	return err
}
