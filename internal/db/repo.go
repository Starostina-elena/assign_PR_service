package db

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) AddUser(name string, isActive bool) error {
	_, err := db.conn.Exec(
		`INSERT INTO users (name, is_active) VALUES ($1, $2)`,
		name, isActive,
	)
	db.log.Info("user added", "name", name, "is_active", isActive)
	return err
}

func (db *DB) GetUserByID(id int64) (string, bool, error) {
	var name string
	var isActive bool
	err := db.conn.QueryRow(
		`SELECT name, is_active FROM users WHERE id = $1`,
		id,
	).Scan(&name, &isActive)
	if err != nil {
		return "", false, err
	}
	return name, isActive, nil
}
