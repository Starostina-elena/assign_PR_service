package core

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

type User struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

type Team struct {
	ID      int64
	Name    string
	Members []User
}

type PullRequest struct {
	ID          int64
	Title       string
	IsOpened    bool
	AuthorID    int64
	Reviewer1ID int64
	Reviewer2ID int64
}
