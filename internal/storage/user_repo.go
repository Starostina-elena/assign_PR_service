package storage

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"context"
)

func (db *DB) AddUser(ctx context.Context, name string, isActive bool) (int64, error) {
	var id int64
	err := db.conn.QueryRowContext(
		ctx,
		`INSERT INTO users (name, is_active) VALUES ($1, $2) RETURNING id`,
		name, isActive,
	).Scan(&id)
	if err != nil {
		db.log.Error("failed to add user", "name", name, "is_active", isActive, "error", err)
		return 0, err
	}
	db.log.Info("user added", "name", name, "is_active", isActive)
	return id, err
}

func (db *DB) GetUserByID(ctx context.Context, id int64) (core.User, error) {
	var u core.User
	err := db.conn.GetContext(ctx, &u, `SELECT id, name, is_active FROM users WHERE id = $1`, id)
	if err != nil {
		return core.User{}, err
	}
	return u, nil
}
