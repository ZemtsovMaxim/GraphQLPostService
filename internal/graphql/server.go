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
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return postService.GetAllPosts(), nil
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
					title := params.Args["title"].(string)
					content := params.Args["content"].(string)
					return postService.CreatePost(title, content), nil
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
					id := params.Args["id"].(int)
					return postService.DisableComments(id), nil
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
