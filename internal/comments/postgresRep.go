package comments

import (
	"database/sql"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

// Реализация методов интерфейса CommentRepository для PostgreSQL

func (r *PostgresCommentRepository) CreateComment(comment *Comment) error {
	_, err := r.db.Exec("INSERT INTO comments (post_id, text) VALUES ($1, $2)", comment.PostID, comment.Text)
	return err
}

func (r *PostgresCommentRepository) GetCommentsByPostID(postID int, limit, offset int) ([]*Comment, error) {
	rows, err := r.db.Query("SELECT id, post_id, text FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3", postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Text); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
