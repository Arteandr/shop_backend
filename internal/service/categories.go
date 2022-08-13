package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
)

type CategoriesService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) Create(ctx context.Context, name string) (int, error) {
	category := models.Category{
		Name: name,
	}
	id, err := s.repo.Create(ctx, category)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *CategoriesService) Delete(ctx context.Context, categoryId int) error {
	return s.repo.Delete(ctx, categoryId)
}

func (s *CategoriesService) GetAll(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoriesService) Exist(ctx context.Context, colorId int) (bool, error) {
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

func (s *CategoriesService) GetById(ctx context.Context, categoryId int) (models.Category, error) {
	var category models.Category
	return category, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		exist, err := s.repo.Exist(ctx, categoryId)
		if !exist {
			return apperrors.ErrIdNotFound("category", categoryId)
		} else if err != nil {
			return err
		}

		category, err = s.repo.GetById(ctx, categoryId)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *CategoriesService) Update(ctx context.Context, categoryId int, name string) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.repo.Exist(ctx, categoryId)
		if !exist {
			return apperrors.ErrIdNotFound("category", categoryId)
		} else if err != nil {
			return err
		}

		category := models.Category{
			Id:   categoryId,
			Name: name,
		}

		return s.repo.Update(ctx, category)
	})
}
