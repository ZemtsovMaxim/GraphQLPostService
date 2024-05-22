package myGraphql

import (
	"context"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
)

type Resolver struct {
	PostService    *posts.PostService
	CommentService *comments.CommentService
}

func NewResolver(postService *posts.PostService, commentService *comments.CommentService) *Resolver {
	return &Resolver{
		PostService:    postService,
		CommentService: commentService,
	}
}

// Резолвер для Query.posts
func (r *Resolver) Posts(ctx context.Context) ([]*posts.Post, error) {
	return r.PostService.GetAllPosts()
}

// Резолвер для Query.commentsByPostID
func (r *Resolver) CommentsByPostID(ctx context.Context, postID int, limit, offset *int) ([]*comments.Comment, error) {
	return r.CommentService.GetCommentsByPostID(postID, *limit, *offset), nil
}

// Резолвер для Query.getPostByID
func (r *Resolver) GetPostByID(ctx context.Context, id int) (*posts.Post, error) {
	return r.PostService.GetPostByID(id)
}

// Резолвер для Mutation.createPost
func (r *Resolver) CreatePost(ctx context.Context, title string, content string) error {
	return r.PostService.CreatePost(title, content)
}

// Резолвер для Mutation.createComment
func (r *Resolver) CreateComment(ctx context.Context, postID int, text string) (bool, error) {
	return r.CommentService.CreateComment(postID, text), nil
}

// Резолвер для Mutation.disableComments
func (r *Resolver) DisableComments(ctx context.Context, postID int) (bool, error) {
	err := r.PostService.DisableComments(postID)
	if err != nil {
		return false, err
	}
	return true, nil
}
