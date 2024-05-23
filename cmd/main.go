package main

import (
	"log/slog"
	"net/http"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/comments"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/config"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/database"
	myGraphql "github.com/ZemtsovMaxim/OzonTestTask/internal/graphql"
	"github.com/ZemtsovMaxim/OzonTestTask/internal/logger"
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
		db, err := database.Connect(cfg.DB, log)
		if err != nil {
			log.Error("failed to connect to the database", slog.Any("err", err))
			return
		}
		postRepo = posts.NewPostgresPostRepository(db, log)
		commentRepo = comments.NewPostgresCommentRepository(db, log)
	} else if cfg.Storage == "in-memory" {
		postRepo = posts.NewInMemoryPostRepository(log)
		commentRepo = comments.NewInMemoryCommentRepository(log)
	} else {
		log.Error("invalid storage type", slog.String("storage", cfg.Storage))
		return
	}

	postService := posts.NewPostService(postRepo, log)
	commentService := comments.NewCommentService(commentRepo, log)

	// Настраиваем GraphQL сервер
	server := myGraphql.NewServer(postService, commentService, log)

	// Запускаем HTTP сервер
	http.Handle("/", server)
	log.Info("starting server", slog.String("address", cfg.Address))
	err := http.ListenAndServe(cfg.Address, nil)
	if err != nil {
		log.Error("failed to start server", slog.Any("err", err))
	}
}
