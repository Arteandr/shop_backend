package service

import (
	"context"

	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
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

func (s *ColorsService) Delete(ctx context.Context, colorsId []int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, colorId := range colorsId {
			if err := s.repo.Delete(ctx, colorId); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ColorsService) DeleteFromItems(ctx context.Context, colorId int) error {
	return s.repo.DeleteFromItems(ctx, colorId)
}

func (s *ColorsService) AddToItems(ctx context.Context, colorId int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.repo.Exist(ctx, colorId)
		if !exist {
			return apperrors.ErrIdNotFound("color", colorId)
		} else if err != nil {
			return err
		}

		return s.repo.AddToItems(ctx, colorId)
	})
}

func (s *ColorsService) Update(ctx context.Context, id int, name, hex string, price float64) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.repo.Exist(ctx, id)
		if !exist {
			return apperrors.ErrIdNotFound("color", id)
		} else if err != nil {
			return err
		}

		color := models.Color{
			Id:    id,
			Name:  name,
			Hex:   hex,
			Price: price,
		}

		return s.repo.Update(ctx, color)
	})
}

func (s *ColorsService) Exist(ctx context.Context, colorId int) (bool, error) {
	var exist bool

	return exist, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		exist, err = s.repo.Exist(ctx, colorId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *ColorsService) GetById(ctx context.Context, colorId int) (models.Color, error) {
	var color models.Color

	return color, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		exist, err := s.repo.Exist(ctx, colorId)
		if !exist {
			return apperrors.ErrIdNotFound("color", colorId)
		} else if err != nil {
			return err
		}

		color, err = s.repo.GetById(ctx, colorId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *ColorsService) GetAll(ctx context.Context) ([]models.Color, error) {
	return s.repo.GetAll(ctx)
}
