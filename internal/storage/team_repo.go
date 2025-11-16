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

func (db *DB) GetTeamMembers(ctx context.Context, teamId int64) ([]core.User, error) {
	var users []core.User
	err := db.conn.SelectContext(ctx, &users, `SELECT id, name, is_active, team_id FROM users WHERE team_id = $1`, teamId)
	if err != nil {
		return nil, err
	}
	return users, nil
}
