package myGraphql

import (
	"net/http"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Server struct {
	*http.ServeMux
}

func NewServer(postService *posts.PostService, commentService *comments.CommentService) *Server {
	srv := &Server{
		ServeMux: http.NewServeMux(),
	}

	schema, err := createSchema(postService, commentService)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})

	srv.Handle("/", h)
	return srv
}

func createSchema(postService *posts.PostService, commentService *comments.CommentService) (*graphql.Schema, error) {
	resolver := NewResolver(postService, commentService)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolvePosts(params)
				},
			},
			"post": &graphql.Field{
				Type: graphql.NewNonNull(postType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolvePost(params)
				},
			},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
				Args: graphql.FieldConfigArgument{
					"postID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"offset": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveCommentsByPostID(params)
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createPost": &graphql.Field{
				Type:        postType,
				Description: "Create a new post",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveCreatePost(params)
				},
			},
			"disableComments": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Disable comments for a post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveDisableComments(params)
				},
			},
			"createComment": &graphql.Field{
				Type:        commentType,
				Description: "Create a new comment",
				Args: graphql.FieldConfigArgument{
					"postID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveCreateComment(params)
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
		"commentsDisabled": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"postID": &graphql.Field{
			Type: graphql.Int,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
	},
})
