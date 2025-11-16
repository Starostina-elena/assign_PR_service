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

func (db *DB) AddTeamWithMembers(ctx context.Context, name string, members []core.User) (core.Team, []core.User, error) {
	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		db.log.Error("failed to begin tx for AddTeamWithMembers", "error", err)
		return core.Team{}, nil, err
	}

	var teamID int64
	err = tx.QueryRowContext(ctx, `INSERT INTO teams (name) VALUES ($1) RETURNING id`, name).Scan(&teamID)
	if err != nil {
		tx.Rollback()
		db.log.Error("failed to insert team", "name", name, "error", err)
		return core.Team{}, nil, err
	}

	for _, m := range members {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO users (id, name, is_active, team_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, is_active = EXCLUDED.is_active, team_id = EXCLUDED.team_id
		`, m.ID, m.Name, m.IsActive, teamID)
		if err != nil {
			tx.Rollback()
			db.log.Error("failed to upsert user", "user_id", m.ID, "error", err)
			return core.Team{}, nil, err
		}
	}

	var users []core.User
	err = tx.SelectContext(ctx, &users, `SELECT id, name, is_active, team_id FROM users WHERE team_id = $1`, teamID)
	if err != nil {
		tx.Rollback()
		db.log.Error("failed to select team members", "team_id", teamID, "error", err)
		return core.Team{}, nil, err
	}

	if err := tx.Commit(); err != nil {
		db.log.Error("failed to commit AddTeamWithMembers tx", "error", err)
		return core.Team{}, nil, err
	}

	team := core.Team{ID: teamID, Name: name}
	return team, users, nil
}
