package comments

import (
	"database/sql"
	"log/slog"
)

type PostgresCommentRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresCommentRepository(db *sql.DB, log *slog.Logger) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db, logger: log}
}

// Реализация методов интерфейса CommentRepository для PostgreSQL

func (r *PostgresCommentRepository) CreateComment(comment *Comment) error {
	_, err := r.db.Exec("INSERT INTO comments (post_id, text) VALUES ($1, $2)", comment.PostID, comment.Text)
	if err != nil {
		r.logger.Error("Cant create comment", slog.Any("err", err))
	}
	return err
}

func (r *PostgresCommentRepository) GetCommentsByPostID(postID int, limit, offset int) ([]*Comment, error) {
	rows, err := r.db.Query("SELECT id, post_id, text FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3", postID, limit, offset)
	if err != nil {
		r.logger.Error("Cant get comment by post", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Text); err != nil {
			r.logger.Error("Cant scan comment", slog.Any("err", err))
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
		r.logger.Error("Cant get replies", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var replies []*Comment
	for rows.Next() {
		var reply Comment
		err := rows.Scan(&reply.ID, &reply.PostID, &reply.Text, &reply.ParentID)
		if err != nil {
			r.logger.Error("Cant scan reply", slog.Any("err", err))
			return nil, err
		}
		replies = append(replies, &reply)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Error with rows", slog.Any("err", err))
		return nil, err
	}

	return replies, nil
}

func (r *PostgresCommentRepository) GetCommentByID(id int) (*Comment, error) {
	var comment Comment
	err := r.db.QueryRow("SELECT id, post_id, content, parent_id FROM comments WHERE id = $1", id).Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.ParentID)
	if err == sql.ErrNoRows {
		r.logger.Error("sql: no rows in result set", slog.Any("err", err))
		return nil, err
	}
	if err != nil {
		r.logger.Error("Cant get comment by id", slog.Any("err", err))
		return nil, err
	}
	return &comment, nil
}
