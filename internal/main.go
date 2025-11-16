package internal

import (
	"fmt"
	"log/slog"
	"os"

	storage "Starostina-elena/pull_req_assign/internal/storage"
)

func Init() error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("initializing db connection")

	dbHost := getenv("DB_HOST", "db")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "postgres")
	dbPass := getenv("DB_PASSWORD", "postgres")
	dbName := getenv("DB_NAME", "postgres")

	db_link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := storage.New(logger, db_link)
	if err != nil {
		logger.Error("db connect failed", "error", err)
		return err
	}

	if err := db.Migrate(); err != nil {
		logger.Error("migrate failed", "error", err)
		return err
	}

	logger.Info("db initialized and migrations applied")
	defer db.Close()

	db.AddUser("Vasya", true)
	name, isActive, err := db.GetUserByID(1)
	if err != nil {
		logger.Error("failed to get user by ID", "error", err)
	} else {
		logger.Info("user retrieved", "name", name, "is_active", isActive)
	}
	return nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
