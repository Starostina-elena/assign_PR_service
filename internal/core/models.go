package core

import "database/sql"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

type User struct {
	ID       string        `db:"id" json:"id"`
	Name     string        `db:"name" json:"name"`
	TeamID   sql.NullInt64 `db:"team_id" json:"team_id"`
	IsActive bool          `db:"is_active" json:"is_active"`
}

type Team struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type PullRequest struct {
	ID          int64          `db:"id" json:"id"`
	Title       string         `db:"title" json:"title"`
	IsOpened    bool           `db:"is_opened" json:"is_opened"`
	AuthorID    string         `db:"author_id" json:"author_id"`
	Reviewer1ID sql.NullString `db:"reviewer1_id" json:"reviewer1_id"`
	Reviewer2ID sql.NullString `db:"reviewer2_id" json:"reviewer2_id"`
	CreatedAt   string         `db:"created_at" json:"created_at"`
	MergedAt    string         `db:"merged_at" json:"merged_at"`
}
