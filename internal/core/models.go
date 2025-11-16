package core

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type User struct {
	ID       int64
	Name     string
	IsActive bool
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
