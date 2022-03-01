package service

import (
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
)

type Users interface {
	EmailExist(email string) bool
	SignUp(email, password string) (int, error)
}

type Services struct {
	Users Users
}

type ServicesDeps struct {
	Repos        *repository.Repositories
	TokenManager auth.TokenManager
	Hasher       hash.PasswordHasher
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager),
	}
}
