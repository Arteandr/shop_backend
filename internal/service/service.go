package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"time"
)

type Colors interface {
	Exist(colorId int) (bool, error)
	Create(name, hex string, price float64) (int, error)
}

type Categories interface {
	Exist(colorId int) (bool, error)
	Create(name string) (int, error)
}

type Items interface {
	Create(name, description string, categoryId int, sku string) (int, error)
	LinkColor(itemId int, colorId int) error
	LinkTags(itemId int, tags []string) error
	GetById(itemId int) (models.Item, error)
	GetBySku(sku string) (models.Item, error)
}

type Users interface {
	EmailExist(email string) bool
	SignUp(email, password string) (int, error)
	SignIn(email, password string) (models.Tokens, error)
	GetUserById(id int) (models.User, error)
}

type Services struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
}

type ServicesDeps struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	Hasher          hash.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users:      NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Items:      NewItemsService(deps.Repos.Items),
		Categories: NewCategoriesService(deps.Repos.Categories),
		Colors:     NewColorsService(deps.Repos.Colors),
	}
}
