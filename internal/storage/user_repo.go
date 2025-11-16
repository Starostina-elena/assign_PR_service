package storage

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"context"
)

func (db *DB) AddUser(ctx context.Context, id string, name string, isActive bool) error {
	err := db.conn.QueryRowContext(
		ctx,
		`INSERT INTO users (id, name, is_active) VALUES ($1, $2, $3) RETURNING id`,
		id, name, isActive,
	).Scan(&id)
	if err != nil {
		db.log.Error("failed to add user", "name", name, "is_active", isActive, "error", err)
		return err
	}
	db.log.Info("user added", "name", name, "is_active", isActive)
	return nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (core.User, error) {
	var u core.User
	err := db.conn.GetContext(ctx, &u, `SELECT id, name, team_id, is_active FROM users WHERE id = $1`, id)
	if err != nil {
		return core.User{}, err
	}
	return u, nil
}

func (db *DB) SetTeamToUser(ctx context.Context, userId string, teamId int64) error {
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

func (db *DB) ExpellUserFromTeam(ctx context.Context, userId string) error {
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

func (db *DB) SetUserIsActive(ctx context.Context, userId string, isActive bool) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2`,
		isActive, userId,
	)
	if err != nil {
		db.log.Error("failed to set user is active", "user_id", userId, "error", err)
		return err
	}
	db.log.Info("user is active set", "user_id", userId, "is_active", isActive)
	return nil
}

func (dn *DB) GetActiveCoworkers(ctx context.Context, userId string) ([]core.User, error) {
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

func (db *DB) GetPullRequestsAssigned(ctx context.Context, reviewerId string) ([]core.PullRequest, error) {
	var prs []core.PullRequest
	err := db.conn.SelectContext(ctx, &prs, `
		SELECT id, title, is_opened, author_id, reviewer1_id, reviewer2_id
		FROM pull_requests
		WHERE is_opened = TRUE
		  AND (reviewer1_id = $1 OR reviewer2_id = $1)
	`, reviewerId)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
