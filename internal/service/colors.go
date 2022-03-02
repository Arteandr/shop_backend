package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type ColorsService struct {
	repo repository.Colors
}

func NewColorsService(repo repository.Colors) *ColorsService {
	return &ColorsService{repo: repo}
}

func (s *ColorsService) Create(name, hex string, price float64) (int, error) {
	color := models.Color{
		Name:  name,
		Hex:   hex,
		Price: price,
	}

	id, err := s.repo.Create(color)
	if err != nil {
		return 0, err
	}

	return id, nil
}
