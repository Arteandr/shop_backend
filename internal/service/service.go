package service

import "shop_backend/internal/repository"

type Services struct {
}

type ServicesDeps struct {
	Repos *repository.Repositories
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{}
}
