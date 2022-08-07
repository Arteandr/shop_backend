package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

const (
	usersTable         = "users"
	categoriesTable    = "categories"
	itemsTable         = "items"
	colorsTable        = "colors"
	itemsColorsTable   = "items_colors"
	tagsTable          = "tags"
	imagesTable        = "images"
	itemsImagesTable   = "items_images"
	sessionsTable      = "sessions"
	addressTable       = "address"
	usersInvoiceTable  = "users_invoice"
	usersShippingTable = "users_shipping"
	phonesTable        = "phone_numbers"
)

type Images interface {
	Upload(filename string) (int, error)
	GetAll() ([]models.Image, error)
	GetById(imageId int) (models.Image, error)
	Exist(imageId int) (bool, error)
	Delete(imageId int) error
	DeleteFromItems(imageId int) error
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
	Update(itemId int, name, description string, categoryId int, price float64, sku string) error
	Delete(itemId int) error
	DeleteTags(itemId int) error
	DeleteImages(itemId int) error
	DeleteColors(itemId int) error
	Exist(itemId int) (bool, error)
}

type Users interface {
	SetSession(ctx context.Context, userId int, session models.Session) error
	DeleteSession(ctx context.Context, userId int) error
	Create(ctx context.Context, user models.User) (models.User, error)
	CreatePhone(ctx context.Context, userId int) error
	GetByCredentials(ctx context.Context, findBy, login, password string) (models.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	GetById(ctx context.Context, userId int) (models.User, error)
	GetPhone(ctx context.Context, userId int) (models.Phone, error)
	GetAddress(ctx context.Context, typeof string, userId int) (models.Address, error)
	UpdateField(ctx context.Context, field string, value interface{}, userId int) error
	UpdatePhone(ctx context.Context, phoneCode, phoneNumber string, userId int) error
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
