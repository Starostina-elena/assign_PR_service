package main

import (
	"Starostina-elena/pull_req_assign/internal"
	"Starostina-elena/pull_req_assign/internal/core"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg core.Config) {
	slog.Info("Server started")

	app, err := internal.Init(cfg)
	if err != nil {
		slog.Error("Initialization failed", "error", err)
		return
	}
	app.Logger.Info("Initialization succeeded")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		app.Logger.Error("Run failed", "error", err)
	} else {
		app.Logger.Info("Run exited")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*1e9)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		app.Logger.Error("Stop failed", "error", err)
	} else {
		app.Logger.Info("Stop succeeded")
	}
}

func main() {
	cfg := core.Config{
		DBHost:     getenv("DB_HOST", "db"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBUser:     getenv("DB_USER", "postgres"),
		DBPassword: getenv("DB_PASSWORD", "postgres"),
		DBName:     getenv("DB_NAME", "postgres"),
		AppPort:    getenv("APP_PORT", "8080"),
	}
	Run(cfg)
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
