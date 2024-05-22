package myGraphql

import (
	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
	"github.com/graphql-go/graphql"
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

func (r *Resolver) resolvePosts(p graphql.ResolveParams) (interface{}, error) {
	return r.PostService.GetAllPosts()
}

func (r *Resolver) resolvePost(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
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

func (r *Resolver) resolveCreateComment(p graphql.ResolveParams) (interface{}, error) {
	postID, _ := p.Args["postID"].(int)
	content, _ := p.Args["content"].(string)
	return r.CommentService.CreateComment(postID, content)
}
