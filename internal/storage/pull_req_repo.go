package storage

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"context"
	"database/sql"
)

func (db *DB) AddPullRequest(
	ctx context.Context, title string, authorId int64, isOpened bool, rev1, rev2 *int64,
) (int64, error) {
	var id int64
	var r1, r2 sql.NullInt64
	if rev1 != nil {
		r1 = sql.NullInt64{Int64: *rev1, Valid: true}
	}
	if rev2 != nil {
		r2 = sql.NullInt64{Int64: *rev2, Valid: true}
	}

	err := db.conn.QueryRowContext(
		ctx,
		`INSERT INTO pull_requests (title, author_id, is_opened, reviewer1_id, reviewer2_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, authorId, isOpened, r1, r2,
	).Scan(&id)
	if err != nil {
		db.log.Error("failed to add pull request", "title", title, "author", authorId)
		return 0, err
	}
	db.log.Info("pull request added", "title", title, "author", authorId)
	return id, err
}

func (db *DB) GetPullRequestByID(ctx context.Context, id int64) (core.PullRequest, error) {
	var pr core.PullRequest
	err := db.conn.GetContext(ctx, &pr, `SELECT id, title, is_opened, author_id, reviewer1_id, reviewer2_id FROM pull_requests WHERE id = $1`, id)
	if err != nil {
		return core.PullRequest{}, err
	}
	return pr, nil
}
