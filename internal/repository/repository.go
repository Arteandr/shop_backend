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
	tagsTable       = "tags"
)

type Colors interface {
	Exist(colorId int) (bool, error)
	Create(color models.Color) (int, error)
}

type Categories interface {
	Exist(categoryId int) (bool, error)
	Create(category models.Category) (int, error)
}

type Items interface {
	Create(item models.Item) (int, error)
	LinkColor(itemId int, colorId int) error
	LinkTag(itemId int, tag string) error
	GetById(itemId int) (models.Item, error)
	GetBySku(sku string) (models.Item, error)
	GetColors(itemId int) ([]models.Color, error)
	GetTags(itemId int) ([]models.Tag, error)
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
