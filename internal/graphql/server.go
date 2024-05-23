package myGraphql

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/config"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/logger"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Server struct {
	*http.ServeMux
	upgrader websocket.Upgrader
	mu       sync.Mutex
	clients  map[*websocket.Conn]bool
	logger   *slog.Logger
}

// NewServer создает новый экземпляр сервера GraphQL.
func NewServer(postService *posts.PostService, commentService *comments.CommentService, log *slog.Logger) *Server {
	srv := &Server{
		logger:   log,
		ServeMux: http.NewServeMux(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Создаем схему GraphQL
	schema, err := createSchema(postService, commentService)
	if err != nil {
		log.Error("cant create schema", slog.Any("err", err))
	}

	// Добавляем обработчик GraphQL API
	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})

	srv.Handle("/", h)
	srv.HandleFunc("/subscriptions", srv.handleSubscriptions)
	return srv
}

func (srv *Server) handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	conn, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}

	srv.mu.Lock()
	srv.clients[conn] = true
	srv.mu.Unlock()

	defer func() {
		srv.mu.Lock()
		delete(srv.clients, conn)
		srv.mu.Unlock()
		conn.Close()
	}()

	go srv.processSubscriptions(conn)
}

func (srv *Server) processSubscriptions(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			srv.logger.Error("error reading message", slog.Any("err", err))
			return
		}

		srv.logger.Info("Received message:", slog.Any("msg", msg))

		response := map[string]interface{}{
			"type": "data",
			"id":   "1",
			"payload": map[string]interface{}{
				"data": map[string]interface{}{
					"commentAdded": map[string]interface{}{
						"id":       1,
						"text":     "This is a test comment",
						"parentID": nil,
					},
				},
			},
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			srv.logger.Error("error marshaling response", slog.Any("err", err))
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			srv.logger.Error("error writing message", slog.Any("err", err))
			return
		}
	}
}

func createSchema(postService *posts.PostService, commentService *comments.CommentService) (*graphql.Schema, error) {

	cfg := config.MustLoad()

	log := logger.SetUpLogger(cfg.LogLevel)

	resolver := NewResolver(postService, commentService, log)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolvePosts()
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
			"replies": &graphql.Field{
				Type: graphql.NewList(commentType),
				Args: graphql.FieldConfigArgument{
					"parentID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveReplies(params)
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
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"parentID": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.resolveCreateComment(params)
				},
			},
		},
	})

	rootSubscription := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"commentAdded": &graphql.Field{
				Type:        commentType,
				Description: "Subscribe to new comments added to a specific post",
				Args: graphql.FieldConfigArgument{
					"postID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return resolver.postCommentAdded(params)
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        rootQuery,
		Mutation:     rootMutation,
		Subscription: rootSubscription,
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
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"parentID": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
