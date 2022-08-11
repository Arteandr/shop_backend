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

func (s *ItemsService) Create(ctx context.Context, name, description string, categoryId int, sku string, price float64) (int, error) {
	item := models.Item{
		Name:        name,
		Description: description,
		Category:    models.Category{Id: categoryId},
		Price:       price,
		Sku:         sku,
	}

	id, err := s.repo.Create(ctx, item)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *ItemsService) LinkColor(ctx context.Context, itemId int, colorId int) error {
	return s.repo.LinkColor(ctx, itemId, colorId)
}

func (s *ItemsService) LinkImages(ctx context.Context, itemId int, imagesId []int) error {
	for _, imageId := range imagesId {
		if err := s.repo.LinkImage(ctx, itemId, imageId); err != nil {
			return err
		}
	}

	return nil
}

func (s *ItemsService) LinkTags(ctx context.Context, itemId int, tags []string) error {
	for _, tag := range tags {
		if err := s.repo.LinkTag(ctx, itemId, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *ItemsService) GetAll(ctx context.Context, sortOptions models.SortOptions) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetAll(ctx, sortOptions)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}

		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
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

func (s *ItemsService) GetNew(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetNew(ctx, 4)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}

		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
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

func (s *ItemsService) GetById(ctx context.Context, itemId int) (models.Item, error) {
	item, err := s.repo.GetById(ctx, itemId)
	if err != nil {
		return models.Item{}, err
	}

	colors, err := s.repo.GetColors(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Colors = colors

	tags, err := s.repo.GetTags(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Tags = tags

	images, err := s.repo.GetImages(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	for i := range images {
		images[i].Filename = "/files/" + images[i].Filename
	}
	item.Images = images

	return item, nil
}

func (s *ItemsService) GetBySku(ctx context.Context, sku string) (models.Item, error) {
	item, err := s.repo.GetBySku(ctx, sku)
	if err != nil {
		return models.Item{}, err
	}

	colors, err := s.repo.GetColors(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Colors = colors

	tags, err := s.repo.GetTags(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	item.Tags = tags

	images, err := s.repo.GetImages(ctx, item.Id)
	if err != nil {
		return models.Item{}, err
	}
	for i := range images {
		images[i].Filename = "/files/" + images[i].Filename
	}
	item.Images = images

	return item, nil
}

func (s *ItemsService) GetByCategory(ctx context.Context, categoryId int) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetByCategory(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
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

func (s *ItemsService) GetByTag(ctx context.Context, tag string) ([]models.Item, error) {
	var items []models.Item
	ids, err := s.repo.GetByTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		item, err := s.repo.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
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

func (s *ItemsService) Update(ctx context.Context, id int, name, description string, categoryId int, tags []string, colorsId []int, price float64, sku string, imagesId []int) error {
	if err := s.repo.Update(ctx, id, name, description, categoryId, price, sku); err != nil {
		return err
	}

	// Update tags
	if err := s.repo.DeleteTags(ctx, id); err != nil {
		return err
	}
	if len(tags) > 0 {
		for _, tag := range tags {
			if err := s.repo.LinkTag(ctx, id, tag); err != nil {
				return err
			}
		}
	}

	// Update colors
	if err := s.repo.DeleteColors(ctx, id); err != nil {
		return err
	}
	for _, colorId := range colorsId {
		if err := s.repo.LinkColor(ctx, id, colorId); err != nil {
			return err
		}
	}

	// Update images
	if err := s.repo.DeleteImages(ctx, id); err != nil {
		return err
	}
	for _, imageId := range imagesId {
		if err := s.repo.LinkImage(ctx, id, imageId); err != nil {
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

func (s *ItemsService) Exist(ctx context.Context, itemId int) (bool, error) {
	return s.repo.Exist(ctx, itemId)
}
