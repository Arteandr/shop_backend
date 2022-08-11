package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
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

func (s *CategoriesService) Exist(ctx context.Context, categoryId int) (bool, error) {
	return s.repo.Exist(ctx, categoryId)
}

func (s *CategoriesService) Delete(ctx context.Context, categoryId int) error {
	return s.repo.Delete(ctx, categoryId)
}

func (s *CategoriesService) GetAll(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAllC(ctx)
}

func (s *CategoriesService) GetById(ctx context.Context, categoryId int) (models.Category, error) {
	return s.repo.GetById(ctx, categoryId)
}

func (s *CategoriesService) Update(ctx context.Context, categoryId int, name string) error {
	category := models.Category{
		Id:   categoryId,
		Name: name,
	}

	return s.repo.Update(ctx, category)
}
