package service

import (
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
)

type Services struct {
}

type ServicesDeps struct {
	Repos        *repository.Repositories
	TokenManager auth.TokenManager
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{}
}
