package service

import (
	"fmt"
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
		Category:    models.Category{Id: categoryId},
		Price:       price,
		Sku:         sku,
	}

	id, err := s.repo.Create(item)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *ItemsService) GetNew() ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetNew(4)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *ItemsService) LinkColor(itemId int, colorId int) error {
	return s.repo.LinkColor(itemId, colorId)
}

func (s *ItemsService) LinkImages(itemId int, imagesId []int) error {
	for _, imageId := range imagesId {
		if err := s.repo.LinkImage(itemId, imageId); err != nil {
			return err
		}
	}

	return nil
}

func (s *ItemsService) LinkTags(itemId int, tags []string) error {
	for _, tag := range tags {
		if err := s.repo.LinkTag(itemId, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *ItemsService) GetById(itemId int) (models.Item, error) {
	item, err := s.repo.GetById(itemId)
	if err != nil {
		fmt.Println("err", err.Error())
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

	images, err := s.repo.GetImages(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	for i := range images {
		images[i].Filename = "/files/" + images[i].Filename
	}
	item.Images = images

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

	images, err := s.repo.GetImages(item.Id)
	if err != nil {
		return models.Item{}, err
	}
	for i := range images {
		images[i].Filename = "/files/" + images[i].Filename
	}
	item.Images = images

	return item, nil
}

func (s *ItemsService) GetByCategory(categoryId int) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetByCategory(categoryId)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(id)
		if err != nil {
			return nil, err
		}
		colors, err := s.repo.GetColors(item.Id)
		if err != nil {
			return nil, err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(item.Id)
		if err != nil {
			return nil, err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(item.Id)
		if err != nil {
			return nil, err
		}
		for i := range images {
			images[i].Filename = "/files/" + images[i].Filename
		}
		item.Images = images

		items = append(items, item)
	}

	return items, err
}

func (s *ItemsService) GetByTag(tag string) ([]models.Item, error) {
	items, err := s.repo.GetByTag(tag)
	if err != nil {
		return []models.Item{}, err
	}

	for i, item := range items {
		colors, err := s.repo.GetColors(item.Id)
		if err != nil {
			return []models.Item{}, err
		}
		tags, err := s.repo.GetTags(item.Id)
		if err != nil {
			return []models.Item{}, err
		}
		items[i].Colors = colors
		items[i].Tags = tags
	}

	return items, nil
}

func (s *ItemsService) Delete(itemId int) error {
	return s.repo.Delete(itemId)
}

func (s *ItemsService) Exist(itemId int) (bool, error) {
	return s.repo.Exist(itemId)
}
