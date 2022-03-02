package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"time"
)

type Users interface {
	EmailExist(email string) bool
	SignUp(email, password string) (int, error)
	SignIn(email, password string) (models.Tokens, error)
	GetUserById(id int) (models.User, error)
}

type Services struct {
	Users Users
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
		Users: NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
