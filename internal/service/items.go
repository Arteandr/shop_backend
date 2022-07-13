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

func (s *ItemsService) Create(name, description string, categoryId int, sku string, price float64) (int, error) {
	item := models.Item{
		Name:        name,
		Description: description,
		CategoryId:  categoryId,
		Price:       price,
		Sku:         sku,
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

func (s *ItemsService) GetById(itemId int) (models.Item, error) {
	item, err := s.repo.GetById(itemId)
	if err != nil {
		return models.Item{}, err
	}

	colors, err := s.repo.GetColors(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Colors = colors

	tags, err := s.repo.GetTags(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Tags = tags

	return item, nil
}

func (s *ItemsService) GetBySku(sku string) (models.Item, error) {
	item, err := s.repo.GetBySku(sku)
	if err != nil {
		return models.Item{}, err
	}

	colors, err := s.repo.GetColors(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Colors = colors

	tags, err := s.repo.GetTags(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Tags = tags

	return item, nil
}

func (s *ItemsService) LinkTags(itemId int, tags []string) error {
	for _, tag := range tags {
		if err := s.repo.LinkTag(itemId, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *ItemsService) Delete(itemId int) error {
	return s.repo.Delete(itemId)
}
