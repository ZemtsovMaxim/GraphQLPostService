package main

import (
	"log/slog"
	"net/http"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/config"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/database"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/logger"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/myGraphql"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/posts"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.MustLoad()

	// Настраиваем логгер
	log := logger.SetUpLogger(cfg.LogLevel)

	// Инициализируем репозитории и сервисы
	var postRepo posts.PostRepository
	var commentRepo comments.CommentRepository

	if cfg.Storage == "postgres" {
		db, err := database.Connect(cfg.DB)
		if err != nil {
			log.Error("failed to connect to the database", slog.String("error", err.Error()))
			return
		}
		postRepo = posts.NewPostgresPostRepository(db)
		commentRepo = comments.NewPostgresCommentRepository(db)
	} else if cfg.Storage == "in-memory" {
		postRepo = posts.NewInMemoryPostRepository()
		commentRepo = comments.NewInMemoryCommentRepository()
	} else {
		log.Error("invalid storage type", slog.String("storage", cfg.Storage))
		return
	}

	postService := posts.NewPostService(postRepo)
	commentService := comments.NewCommentService(commentRepo)

	// Создаем резолвер
	resolver := myGraphql.NewResolver(postService, commentService)

	// Создаем схему GraphQL
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{}), // Заполните это в соответствии с вашей схемой
		Mutation: graphql.NewObject(graphql.ObjectConfig{}), // Заполните это в соответствии с вашей схемой
	})
	if err != nil {
		log.Error("failed to create GraphQL schema", slog.String("error", err.Error()))
		return
	}

	// Настраиваем GraphQL сервер
	server := graphql.NewServer(&schema)

	// Запускаем HTTP сервер
	http.Handle("/", server)
	log.Info("starting server", slog.String("address", cfg.Address))
	err = http.ListenAndServe(cfg.Address, nil)
	if err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
