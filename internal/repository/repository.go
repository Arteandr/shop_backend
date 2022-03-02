package repository

import (
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

const (
	usersTable      = "users"
	categoriesTable = "categories"
	itemsTable      = "items"
	colorsTable     = "colors"
	itemColorsTable = "item_colors"
)

type Colors interface {
	Create(category models.Color) (int, error)
}

type Categories interface {
	Create(category models.Category) (int, error)
}

type Items interface {
}

type Users interface {
	Create(user models.User) (int, error)
	GetByCredentials(email, passwordHash string) (models.User, error)
	GetById(id int) (models.User, error)
	Exist(email string) bool
}

type Repositories struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:      NewUsersRepo(db),
		Items:      NewItemsRepo(db),
		Categories: NewCategoriesRepo(db),
		Colors:     NewColorsRepo(db),
	}
}
