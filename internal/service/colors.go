package service

import "shop_backend/internal/repository"

type ColorsService struct {
	repo repository.Colors
}

func NewColorsService(repo repository.Colors) *ColorsService {
	return &ColorsService{repo: repo}
}
