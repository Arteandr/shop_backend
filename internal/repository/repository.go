package repository

import (
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

const (
	usersTable       = "users"
	categoriesTable  = "categories"
	itemsTable       = "items"
	colorsTable      = "colors"
	itemsColorsTable = "items_colors"
	tagsTable        = "tags"
	imagesTable      = "images"
	itemsImagesTable = "items_images"
)

type Images interface {
	Upload(filename string) (int, error)
	GetAll() ([]models.Image, error)
	Exist(imageId int) (bool, error)
}

type Colors interface {
	Exist(colorId int) (bool, error)
	GetById(colorId int) (models.Color, error)
	GetAll() ([]models.Color, error)
	Create(color models.Color) (int, error)
	Update(color models.Color) error
	Delete(colorId int) error
	DeleteFromItems(colorId int) error
	AddToItems(colorId int) error
}

type Categories interface {
	Exist(categoryId int) (bool, error)
	Create(category models.Category) (int, error)
	GetAll() ([]models.Category, error)
	Delete(categoryId int) error
	GetById(categoryId int) (models.Category, error)
}

type Items interface {
	Create(item models.Item) (int, error)
	LinkColor(itemId, colorId int) error
	LinkTag(itemId int, tag string) error
	LinkImage(itemId, imageId int) error
	GetNew(limit int) ([]int, error)
	GetById(itemId int) (models.Item, error)
	GetBySku(sku string) (models.Item, error)
	GetByCategory(categoryId int) ([]int, error)
	GetByTag(tag string) ([]int, error)
	GetColors(itemId int) ([]models.Color, error)
	GetTags(itemId int) ([]models.Tag, error)
	GetImages(itemId int) ([]models.Image, error)
	Delete(itemId int) error
	Exist(itemId int) (bool, error)
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
	Images     Images
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:      NewUsersRepo(db),
		Items:      NewItemsRepo(db),
		Categories: NewCategoriesRepo(db),
		Colors:     NewColorsRepo(db),
		Images:     NewImagesRepo(db),
	}
}
