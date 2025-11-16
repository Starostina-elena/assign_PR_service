package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"Starostina-elena/pull_req_assign/internal/api"
	"Starostina-elena/pull_req_assign/internal/core"
	service "Starostina-elena/pull_req_assign/internal/service"
	storage "Starostina-elena/pull_req_assign/internal/storage"
)

type App struct {
	DB     *storage.DB
	Logger *slog.Logger
	Server *http.Server
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

	userService := service.NewUserService(logger, db)
	teamService := service.NewTeamService(logger, db)

	r := api.NewRouter(logger, userService, teamService)

	server := &http.Server{Addr: ":" + cfg.AppPort, Handler: r}

	return &App{DB: db, Logger: logger, Server: server}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("application starting")

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("server failed", "error", err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if a.Server != nil {
		if err := a.Server.Shutdown(ctx); err != nil {
			a.Logger.Error("server shutdown error", "error", err)
		}
	}
	a.DB.Close()
	a.Logger.Info("application stopped")
	return nil
}
