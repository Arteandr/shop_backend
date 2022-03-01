package repository

import (
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

const (
	usersTable      = "users"
	categoriesTable = "categories"
	itemsTable      = "items"
)

type Users interface {
	Create(user models.User) (int, error)
	Exist(email string) bool
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
