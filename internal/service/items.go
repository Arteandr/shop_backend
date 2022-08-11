package service

import (
	"context"
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

func (s *ItemsService) GetAll(sortOptions models.SortOptions) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetAll(sortOptions)
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

	return items, nil
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

	return items, nil
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
	var items []models.Item
	ids, err := s.repo.GetByTag(tag)
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

func (s *ItemsService) Update(id int, name, description string, categoryId int, tags []string, colorsId []int, price float64, sku string, imagesId []int) error {
	if err := s.repo.Update(id, name, description, categoryId, price, sku); err != nil {
		return err
	}

	// Update tags
	if err := s.repo.DeleteTags(id); err != nil {
		return err
	}
	if len(tags) > 0 {
		for _, tag := range tags {
			if err := s.repo.LinkTag(id, tag); err != nil {
				return err
			}
		}
	}

	// Update colors
	if err := s.repo.DeleteColors(id); err != nil {
		return err
	}
	for _, colorId := range colorsId {
		if err := s.repo.LinkColor(id, colorId); err != nil {
			return err
		}
	}

	// Update images
	if err := s.repo.DeleteImages(id); err != nil {
		return err
	}
	for _, imageId := range imagesId {
		if err := s.repo.LinkImage(id, imageId); err != nil {
			return err
		}
	}

	return nil
}

func (s *ItemsService) Delete(ctx context.Context, itemsId []int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, itemId := range itemsId {
			if err := s.repo.Delete(ctx, itemId); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ItemsService) Exist(itemId int) (bool, error) {
	return s.repo.Exist(itemId)
}
