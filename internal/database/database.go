package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ZemtsovMaxim/OzonTestTask/internal/config"
	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig, logger *slog.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("cant open DB", slog.Any("err", err))
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Error("cant ping DB", slog.Any("err", err))
		return nil, err
	}
	return db, nil
}
