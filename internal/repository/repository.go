package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

const (
	usersTable           = "users"
	categoriesTable      = "categories"
	itemsTable           = "items"
	colorsTable          = "colors"
	itemsColorsTable     = "items_colors"
	tagsTable            = "tags"
	imagesTable          = "images"
	itemsImagesTable     = "items_images"
	sessionsTable        = "sessions"
	addressTable         = "address"
	phonesTable          = "phone_numbers"
	deliveryTable        = "delivery"
	deliveryCompanyTable = "delivery_company"
)

type Images interface {
	Upload(ctx context.Context, filename string) error
	GetAll(ctx context.Context) ([]models.Image, error)
	GetById(ctx context.Context, imageId int) (models.Image, error)
	Exist(ctx context.Context, imageId int) (bool, error)
	Delete(ctx context.Context, imageId int) error
	DeleteFromItems(ctx context.Context, imageId int) error
	Transactor
}

type Colors interface {
	Exist(ctx context.Context, colorId int) (bool, error)
	GetById(ctx context.Context, colorId int) (models.Color, error)
	GetAll(ctx context.Context) ([]models.Color, error)
	Create(ctx context.Context, color models.Color) (int, error)
	Update(ctx context.Context, color models.Color) error
	Delete(ctx context.Context, colorId int) error
	DeleteFromItems(ctx context.Context, colorId int) error
	AddToItems(ctx context.Context, colorId int) error
	Transactor
}

type Categories interface {
	Exist(ctx context.Context, categoryId int) (bool, error)
	Create(ctx context.Context, category models.Category) (int, error)
	GetAllC(ctx context.Context) ([]models.Category, error)
	Delete(ctx context.Context, categoryId int) error
	GetById(ctx context.Context, categoryId int) (models.Category, error)
	Update(ctx context.Context, category models.Category) error
	Transactor
}

type Items interface {
	Create(ctx context.Context, item models.Item) (int, error)
	LinkColor(ctx context.Context, itemId int, colorId int) error
	LinkTag(ctx context.Context, itemId int, tag string) error
	LinkImage(ctx context.Context, itemId, imageId int) error
	GetNew(ctx context.Context, limit int) ([]int, error)
	GetAll(ctx context.Context, sortOptions models.SortOptions) ([]int, error)
	GetById(ctx context.Context, itemId int) (models.Item, error)
	GetBySku(ctx context.Context, sku string) (models.Item, error)
	GetByCategory(ctx context.Context, categoryId int) ([]int, error)
	GetByTag(ctx context.Context, tag string) ([]int, error)
	GetColors(ctx context.Context, itemId int) ([]models.Color, error)
	GetTags(ctx context.Context, itemId int) ([]models.Tag, error)
	GetImages(ctx context.Context, itemId int) ([]models.Image, error)
	Update(ctx context.Context, itemId int, name, description string, categoryId int, price float64, sku string) error
	Delete(ctx context.Context, itemId int) error
	DeleteTags(ctx context.Context, itemId int) error
	DeleteImages(ctx context.Context, itemId int) error
	DeleteColors(ctx context.Context, itemId int) error
	Exist(ctx context.Context, itemId int) (bool, error)
	Transactor
}

type Users interface {
	SetSession(ctx context.Context, userId int, session models.Session) error
	DeleteSession(ctx context.Context, userId int) error
	Delete(ctx context.Context, userId int) error
	Create(ctx context.Context, user models.User) (models.User, error)
	CreatePhone(ctx context.Context, userId int) error
	CreateDefaultAddress(ctx context.Context, table string, userId int) error
	LinkAddress(ctx context.Context, table string, userId int, addressId int) error
	CreateAddress(ctx context.Context, address models.Address) (models.Address, error)
	GetByCredentials(ctx context.Context, findBy, login, password string) (models.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	GetById(ctx context.Context, userId int) (models.User, error)
	GetPhone(ctx context.Context, userId int) (models.Phone, error)
	GetAddress(ctx context.Context, typeof string, userId int) (models.Address, error)
	GetAll(ctx context.Context) ([]models.User, error)
	UpdateField(ctx context.Context, field string, value interface{}, userId int) error
	UpdatePhone(ctx context.Context, phoneCode, phoneNumber string, userId int) error
	Transactor
}

type Delivery interface {
	Create(ctx context.Context, delivery models.Delivery) (int, error)
	CreateCompany(ctx context.Context, name string) error
	ExistCompany(ctx context.Context, name string) (bool, error)
	GetById(ctx context.Context, deliveryId int) (models.Delivery, error)
	GetAll(ctx context.Context) ([]models.Delivery, error)
	Update(ctx context.Context, delivery models.Delivery) error
	Delete(ctx context.Context, deliveryId int) error
	Transactor
}

type Repositories struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
	Images     Images
	Delivery   Delivery
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:      NewUsersRepo(db),
		Items:      NewItemsRepo(db),
		Categories: NewCategoriesRepo(db),
		Colors:     NewColorsRepo(db),
		Images:     NewImagesRepo(db),
		Delivery:   NewDeliveryRepo(db),
	}
}

type SqlxDB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
