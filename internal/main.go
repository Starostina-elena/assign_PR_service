package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"Starostina-elena/pull_req_assign/internal/core"
	storage "Starostina-elena/pull_req_assign/internal/storage"
)

type App struct {
	DB     *storage.DB
	Logger *slog.Logger
}

func Init(cfg core.Config) (*App, error) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("initializing db connection")

	db_link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := storage.New(logger, db_link)
	if err != nil {
		logger.Error("db connect failed", "error", err)
		return nil, err
	}
	if err := db.Migrate(); err != nil {
		logger.Error("migrate failed", "error", err)
		return nil, err
	}

	logger.Info("db initialized and migrations applied")

	return &App{DB: db, Logger: logger}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("application starting")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.DB.Close()
	a.Logger.Info("application stopped")
	return nil
}
