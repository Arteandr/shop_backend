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

func (s *CategoriesService) Exist(categoryId int) (bool, error) {
	return s.repo.Exist(categoryId)
}

func (s *CategoriesService) Delete(categoryId int) error {
	return s.repo.Delete(categoryId)
}

func (s *CategoriesService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoriesService) GetById(categoryId int) (models.Category, error) {
	return s.repo.GetById(categoryId)
}

func (s *CategoriesService) Update(categoryId int, name string) error {
	category := models.Category{
		Id:   categoryId,
		Name: name,
	}

	return s.repo.Update(category)
}
