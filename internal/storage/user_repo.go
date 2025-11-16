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
	err := db.conn.GetContext(ctx, &u, `SELECT id, name, team_id, is_active FROM users WHERE id = $1`, id)
	if err != nil {
		return core.User{}, err
	}
	return u, nil
}

func (db *DB) SetTeamToUser(ctx context.Context, userId, teamId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE users SET team_id = $1 WHERE id = $2`,
		teamId, userId,
	)
	if err != nil {
		db.log.Error("failed to set team to user", "user_id", userId, "team_id", teamId, "error", err)
		return err
	}
	db.log.Info("team set to user", "user_id", userId, "team_id", teamId)
	return nil
}

func (db *DB) ExpellUserFromTeam(ctx context.Context, userId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE users SET team_id = NULL WHERE id = $1`,
		userId,
	)
	if err != nil {
		db.log.Error("failed to expel user from team", "user_id", userId, "error", err)
		return err
	}
	db.log.Info("user expelled from team", "user_id", userId)
	return nil
}

func (db *DB) ActivateUser(ctx context.Context, userId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE users SET is_active = TRUE WHERE id = $1`,
		userId,
	)
	if err != nil {
		db.log.Error("failed to activate user", "user_id", userId, "error", err)
		return err
	}
	db.log.Info("user activated", "user_id", userId)
	return nil
}

func (db *DB) DeactivateUser(ctx context.Context, userId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE users SET is_active = FALSE WHERE id = $1`,
		userId,
	)
	if err != nil {
		db.log.Error("failed to deactivate user", "user_id", userId, "error", err)
		return err
	}
	db.log.Info("user deactivated", "user_id", userId)
	return nil
}

func (dn *DB) GetActiveCoworkers(ctx context.Context, userId int64) ([]core.User, error) {
	var users []core.User
	err := dn.conn.SelectContext(ctx, &users, `
		SELECT id, name, is_active, team_id
		FROM users
		WHERE team_id = (SELECT team_id FROM users WHERE id = $1)
		  AND is_active = TRUE
		  AND id != $1
	`, userId)
	if err != nil {
		return nil, err
	}
	return users, nil
}
