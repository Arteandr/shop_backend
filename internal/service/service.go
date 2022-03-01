package service

import (
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
)

type Services struct {
}

type ServicesDeps struct {
	Repos        *repository.Repositories
	TokenManager auth.TokenManager
	Hasher       hash.PasswordHasher
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{}
}
