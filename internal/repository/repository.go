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
	GetByCredentials(email, passwordHash string) (models.User, error)
	GetById(id int) (models.User, error)
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
