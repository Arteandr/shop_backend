package service

import (
	"context"
	"mime/multipart"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"time"
)

type Images interface {
	Upload(image *multipart.FileHeader) (int, error)
	GetAll() ([]models.Image, error)
	Exist(imageId int) (bool, error)
	Delete(imageId int) error
}

type Colors interface {
	Exist(colorId int) (bool, error)
	GetById(colorId int) (models.Color, error)
	GetAll() ([]models.Color, error)
	Create(name, hex string, price float64) (int, error)
	Update(id int, name, hex string, price float64) error
	Delete(colorId int) error
	DeleteFromItems(colorId int) error
	AddToItems(colorId int) error
}

type Categories interface {
	Exist(categoryId int) (bool, error)
	GetAll() ([]models.Category, error)
	GetById(categoryId int) (models.Category, error)
	Create(name string) (int, error)
	Delete(categoryId int) error
	Update(categoryId int, name string) error
}

type Items interface {
	Create(name, description string, categoryId int, sku string, price float64) (int, error)
	Update(id int, name, description string, categoryId int, tags []string, colorsId []int, price float64, sku string, imagesId []int) error
	LinkColor(itemId int, colorId int) error
	LinkTags(itemId int, tags []string) error
	LinkImages(itemId int, imagesId []int) error
	GetAll(sortOptions models.SortOptions) ([]models.Item, error)
	GetNew() ([]models.Item, error)
	GetById(itemId int) (models.Item, error)
	GetBySku(sku string) (models.Item, error)
	GetByCategory(categoryId int) ([]models.Item, error)
	GetByTag(tag string) ([]models.Item, error)
	Delete(itemId int) error
	Exist(itemId int) (bool, error)
}

type Users interface {
	SignUp(ctx context.Context, email, login, password string) (models.User, error)
	SignIn(ctx context.Context, findBy, login, password string) (models.Tokens, error)
	Logout(ctx context.Context, userId int) error
	GetMe(ctx context.Context, userId int) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (models.Tokens, error)
	UpdateEmail(ctx context.Context, userId int, email string) error
	UpdatePassword(ctx context.Context, userId int, oldPassword, newPassword string) error
	UpdateInfo(ctx context.Context, userId int, login, firstName, lastName, phoneCode, phoneNumber string) error
	UpdateAddress(ctx context.Context, userId int, different bool, invoiceAddress models.Address, shippingAddress models.Address) error
	DeleteMe(ctx context.Context, userId int) error
}

type Delivery interface {
	Create(ctx context.Context, delivery models.Delivery) (int, error)
	GetById(ctx context.Context, deliveryId int) (models.Delivery, error)
}

type Services struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
	Images     Images
	Delivery   Delivery
}

type ServicesDeps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Items:      NewItemsService(deps.Repos.Items),
		Categories: NewCategoriesService(deps.Repos.Categories),
		Colors:     NewColorsService(deps.Repos.Colors),
		Images:     NewImagesService(deps.Repos.Images),
		Delivery:   NewDeliveryService(deps.Repos.Delivery),
		Users:      NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
