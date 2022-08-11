package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type ColorsService struct {
	repo repository.Colors
}

func NewColorsService(repo repository.Colors) *ColorsService {
	return &ColorsService{repo: repo}
}

func (s *ColorsService) Create(ctx context.Context, name, hex string, price float64) (int, error) {
	color := models.Color{
		Name:  name,
		Hex:   hex,
		Price: price,
	}

	id, err := s.repo.Create(ctx, color)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ColorsService) Exist(ctx context.Context, colorId int) (bool, error) {
	return s.repo.Exist(ctx, colorId)
}

func (s *ColorsService) Delete(ctx context.Context, colorId int) error {
	return s.repo.Delete(ctx, colorId)
}

func (s *ColorsService) DeleteFromItems(ctx context.Context, colorId int) error {
	return s.repo.DeleteFromItems(ctx, colorId)
}

func (s *ColorsService) AddToItems(ctx context.Context, colorId int) error {
	return s.repo.AddToItems(ctx, colorId)
}

func (s *ColorsService) Update(ctx context.Context, id int, name, hex string, price float64) error {
	color := models.Color{
		Id:    id,
		Name:  name,
		Hex:   hex,
		Price: price,
	}

	return s.repo.Update(ctx, color)
}

func (s *ColorsService) GetById(ctx context.Context, colorId int) (models.Color, error) {
	return s.repo.GetById(ctx, colorId)
}

func (s *ColorsService) GetAll(ctx context.Context) ([]models.Color, error) {
	return s.repo.GetAll(ctx)
}
