package comments

import (
	"database/sql"
	"fmt"
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

func (r *PostgresCommentRepository) GetReplies(parentID int) ([]*Comment, error) {
	query := "SELECT id, post_id, content, parent_id FROM comments WHERE parent_id = $1"
	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("error getting replies: %w", err)
	}
	defer rows.Close()

	var replies []*Comment
	for rows.Next() {
		var reply Comment
		err := rows.Scan(&reply.ID, &reply.PostID, &reply.Text, &reply.ParentID)
		if err != nil {
			return nil, fmt.Errorf("error scanning reply: %w", err)
		}
		replies = append(replies, &reply)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %w", err)
	}

	return replies, nil
}

func (r *PostgresCommentRepository) GetCommentByID(id int) (*Comment, error) {
	var comment Comment
	err := r.db.QueryRow("SELECT id, post_id, content, parent_id FROM comments WHERE id = $1", id).Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.ParentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
