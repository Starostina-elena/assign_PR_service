package storage

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"context"
)

func (db *DB) AddTeam(ctx context.Context, name string) (int64, error) {
	var id int64
	err := db.conn.QueryRowContext(
		ctx,
		`INSERT INTO teams (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)
	if err != nil {
		db.log.Error("failed to add team", "name", name, "error", err)
		return 0, err
	}
	db.log.Info("team added", "name", name)
	return id, err
}

func (db *DB) GetTeamByID(ctx context.Context, id int64) (core.Team, error) {
	var t core.Team
	err := db.conn.GetContext(ctx, &t, `SELECT id, name FROM teams WHERE id = $1`, id)
	if err != nil {
		return core.Team{}, err
	}
	return t, nil
}
