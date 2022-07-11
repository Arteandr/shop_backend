package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type ItemsService struct {
	repo repository.Items
}

func NewItemsService(repo repository.Items) *ItemsService {
	return &ItemsService{repo: repo}
}

func (s *ItemsService) Create(name, description string, categoryId int, tags []string) (int, error) {
	item := models.Item{
		Name:        name,
		Description: description,
		CategoryId:  categoryId,
		Tags:        tags,
	}

	id, err := s.repo.Create(item)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *ItemsService) LinkColor(itemId int, colorId int) error {
	return s.repo.LinkColor(itemId, colorId)
}
