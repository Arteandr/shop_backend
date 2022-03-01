package repository

import "github.com/jmoiron/sqlx"

type Users interface {
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
