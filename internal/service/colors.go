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

func (s *ColorsService) Exist(colorId int) (bool, error) {
	return s.repo.Exist(colorId)
}

func (s *ColorsService) Delete(colorId int) error {
	return s.repo.Delete(colorId)
}

func (s *ColorsService) DeleteFromItems(colorId int) error {
	return s.repo.DeleteFromItems(colorId)
}

func (s *ColorsService) AddToItems(colorId int) error {
	return s.repo.AddToItems(colorId)
}
