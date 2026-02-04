package identity

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Salt      string
	IsVerified bool
	Created   time.Time
	LastLogin time.Time
	Banned    *time.Time
}

type Token struct {
	ID           string
	UserID       int64
	Value        string
	RefreshValue string
	Issued       time.Time
}

type Verification struct {
	ID     string
	UserID int64
}
