package myGraphql

import (
	"errors"
	"log/slog"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
	"github.com/graphql-go/graphql"
)

type Resolver struct {
	PostService    *posts.PostService
	CommentService *comments.CommentService
	logger         *slog.Logger
}

func NewResolver(postService *posts.PostService, commentService *comments.CommentService, log *slog.Logger) *Resolver {
	return &Resolver{
		PostService:    postService,
		CommentService: commentService,
		logger:         log,
	}
}

// Функция для отправки обновлений подписчикам
func (r *Resolver) postCommentAdded(p graphql.ResolveParams) (interface{}, error) {
	// Получаем параметры подписки
	postID := p.Args["postID"].(int)

	// Получаем последний комментарий к посту
	comment, err := r.CommentService.GetLastCommentForPost(postID)
	if err != nil {
		r.logger.Error("cant get last comment for post", slog.Any("err", err))
		return nil, err
	}

	// Отправляем новый комментарий подписчикам
	return comment, nil
}

func (r *Resolver) resolvePosts() (interface{}, error) {
	return r.PostService.GetAllPosts()
}

func (r *Resolver) resolvePost(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		r.logger.Info("graphql params is not ok")
		return nil, nil
	}
	return r.PostService.GetPostByID(id)
}

func (r *Resolver) resolveCreatePost(p graphql.ResolveParams) (interface{}, error) {
	title, _ := p.Args["title"].(string)
	content, _ := p.Args["content"].(string)
	return r.PostService.CreatePost(title, content)
}

func (r *Resolver) resolveDisableComments(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	return r.PostService.DisableComments(id)
}

func (r *Resolver) resolveCommentsByPostID(p graphql.ResolveParams) (interface{}, error) {
	postID, _ := p.Args["postID"].(int)
	limit, _ := p.Args["limit"].(int)
	offset, _ := p.Args["offset"].(int)
	return r.CommentService.GetCommentsByPostID(postID, limit, offset)
}

func (r *Resolver) resolveCreateComment(params graphql.ResolveParams) (interface{}, error) {
	postID, ok := params.Args["postID"].(int)
	if !ok {
		r.logger.Info("missing or invalid postID argument")
		return nil, errors.New("missing or invalid postID argument")
	}
	post, err := r.PostService.GetPostByID(postID)
	if err != nil {
		r.logger.Error("cant GetPostByID", slog.Any("err", err))
		return nil, err
	}
	if post.CommentsDisabled {
		r.logger.Info("comments are disabled for this post")
		return nil, errors.New("comments are disabled for this post")
	}
	content, ok := params.Args["text"].(string)
	if !ok {
		r.logger.Info("missing or invalid content argument")
		return nil, errors.New("missing or invalid content argument")
	}
	parentID, _ := params.Args["parentID"].(int)

	return r.CommentService.CreateComment(postID, content, &parentID)
}

func (r *Resolver) resolveReplies(params graphql.ResolveParams) (interface{}, error) {
	parentID, ok := params.Args["parentID"].(int)
	if !ok {
		r.logger.Info("missing or invalid parentID argument")
		return nil, errors.New("missing or invalid parentID argument")
	}
	return r.CommentService.GetReplies(parentID)
}
