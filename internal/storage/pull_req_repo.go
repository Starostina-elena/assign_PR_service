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

func (db *DB) MergePullRequest(ctx context.Context, id int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE pull_requests SET is_opened = FALSE WHERE id = $1`,
		id,
	)
	if err != nil {
		db.log.Error("failed to merge pull request", "id", id, "error", err)
		return err
	}
	db.log.Info("pull request merged", "id", id)
	return nil
}

func (db *DB) ChangeReviewer1(ctx context.Context, pullReqId int64, newReviewerId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE pull_requests SET reviewer1_id = $1 WHERE id = $2`,
		newReviewerId, pullReqId,
	)
	if err != nil {
		db.log.Error("failed to reset reviewer1", "pull_request_id", pullReqId, "new_reviewer_id", newReviewerId, "error", err)
		return err
	}
	db.log.Info("reviewer1 reset", "pull_request_id", pullReqId, "new_reviewer_id", newReviewerId)
	return nil
}

func (db *DB) ChangeReviewer2(ctx context.Context, pullReqId int64, newReviewerId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE pull_requests SET reviewer2_id = $1 WHERE id = $2`,
		newReviewerId, pullReqId,
	)
	if err != nil {
		db.log.Error("failed to reset reviewer2", "pull_request_id", pullReqId, "new_reviewer_id", newReviewerId, "error", err)
		return err
	}
	db.log.Info("reviewer2 reset", "pull_request_id", pullReqId, "new_reviewer_id", newReviewerId)
	return nil
}

func (db *DB) ResetReviewer1(ctx context.Context, pullReqId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE pull_requests SET reviewer1_id = NULL WHERE id = $1`,
		pullReqId,
	)
	if err != nil {
		db.log.Error("failed to reset reviewer1", "pull_request_id", pullReqId, "error", err)
		return err
	}
	db.log.Info("reviewer1 reset", "pull_request_id", pullReqId)
	return nil
}

func (db *DB) ResetReviewer2(ctx context.Context, pullReqId int64) error {
	_, err := db.conn.ExecContext(
		ctx,
		`UPDATE pull_requests SET reviewer2_id = NULL WHERE id = $1`,
		pullReqId,
	)
	if err != nil {
		db.log.Error("failed to reset reviewer2", "pull_request_id", pullReqId, "error", err)
		return err
	}
	db.log.Info("reviewer2 reset", "pull_request_id", pullReqId)
	return nil
}
