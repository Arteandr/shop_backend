package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type CategoriesService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) Create(name string) (int, error) {
	category := models.Category{
		Name: name,
	}
	id, err := s.repo.Create(category)
	if err != nil {
		return 0, err
	}

	return id, nil
}
