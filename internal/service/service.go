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
	Create(name string) (int, error)
	Delete(categoryId int) error
	GetById(categoryId int) (models.Category, error)
}

type Items interface {
	Create(name, description string, categoryId int, sku string, price float64) (int, error)
	Update(id int, name, description string, categoryId int, tags []string, colorsId []int, price float64, sku string, imagesId []int) error
	LinkColor(itemId int, colorId int) error
	LinkTags(itemId int, tags []string) error
	LinkImages(itemId int, imagesId []int) error
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
	GetById(ctx context.Context, userId int) (models.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (models.Tokens, error)
}

type Services struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
	Images     Images
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
		Users:      NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
