package service

import (
	"context"

	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
)

type ItemsService struct {
	repo              repository.Items
	categoriesService Categories
	colorsService     Colors
	imagesService     Images
}

func NewItemsService(repo repository.Items, categories Categories, colors Colors, images Images) *ItemsService {
	return &ItemsService{
		repo:              repo,
		categoriesService: categories,
		colorsService:     colors,
		imagesService:     images,
	}
}

func (s *ItemsService) Create(ctx context.Context, item models.Item) (models.Item, error) {
	var newItem models.Item
	return newItem, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check category exist
		exist, err := s.categoriesService.Exist(ctx, item.Category.Id)
		if !exist {
			return apperrors.ErrIdNotFound("category", item.Category.Id)
		} else if err != nil {
			return err
		}

		// Check colors exist
		for _, color := range item.Colors {
			exist, err := s.colorsService.Exist(ctx, color.Id)
			if !exist {
				return apperrors.ErrIdNotFound("color", item.Category.Id)
			} else if err != nil {
				return err
			}
		}

		// Check images exist
		for _, image := range item.Images {
			exist, err := s.imagesService.Exist(ctx, image.Id)
			if !exist {
				return apperrors.ErrIdNotFound("image", item.Category.Id)
			} else if err != nil {
				return err
			}
		}

		newUserId, err := s.repo.Create(ctx, item)
		if err != nil {
			return err
		}

		// Link colors
		if err := s.LinkColors(ctx, newUserId, item.Colors); err != nil {
			return err
		}

		// Link tags if more than zero
		if len(item.Tags) > 0 {
			if err := s.LinkTags(ctx, newUserId, item.Tags); err != nil {
				return err
			}
		}

		// Link images
		if err := s.LinkImages(ctx, newUserId, item.Images); err != nil {
			return err
		}

		// Get item
		newItem, err = s.GetById(ctx, newUserId)
		if err != nil {
			return err
		}

		category, err := s.categoriesService.GetById(ctx, item.Category.Id)
		if err != nil {
			return err
		}

		newItem.Category = category

		return nil
	})
}

func (s *ItemsService) LinkColors(ctx context.Context, itemId int, colors []models.Color) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, color := range colors {
			if err := s.repo.LinkColor(ctx, itemId, color.Id); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ItemsService) LinkImages(ctx context.Context, itemId int, images []models.Image) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, image := range images {
			if err := s.repo.LinkImage(ctx, itemId, image.Id); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ItemsService) LinkTags(ctx context.Context, itemId int, tags []models.Tag) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, tag := range tags {
			if err := s.repo.LinkTag(ctx, itemId, tag.Name); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ItemsService) GetAll(ctx context.Context, sortOptions models.SortOptions) ([]models.Item, error) {
	var items []models.Item
	return items, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		ids, err := s.repo.GetAll(ctx, sortOptions)
		if err != nil {
			return err
		}

		for _, id := range ids {
			item, err := s.repo.GetById(ctx, id)
			if err != nil {
				return err
			}

			colors, err := s.repo.GetColors(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Colors = colors

			tags, err := s.repo.GetTags(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Tags = tags

			images, err := s.repo.GetImages(ctx, item.Id)
			if err != nil {
				return err
			}
			for i := range images {
				images[i].Filename = "/files/" + images[i].Filename
			}
			item.Images = images

			category, err := s.categoriesService.GetById(ctx, item.Category.Id)
			if err != nil {
				return err
			}
			item.Category = category

			items = append(items, item)
		}

		return nil
	})
}

func (s *ItemsService) GetNew(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	return items, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		ids, err := s.repo.GetNew(ctx, 4)
		if err != nil {
			return err
		}

		for _, id := range ids {
			item, err := s.repo.GetById(ctx, id)
			if err != nil {
				return err
			}

			colors, err := s.repo.GetColors(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Colors = colors

			tags, err := s.repo.GetTags(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Tags = tags

			images, err := s.repo.GetImages(ctx, item.Id)
			if err != nil {
				return err
			}
			for i := range images {
				images[i].Filename = "/files/" + images[i].Filename
			}
			item.Images = images

			category, err := s.categoriesService.GetById(ctx, item.Category.Id)
			if err != nil {
				return err
			}
			item.Category = category

			items = append(items, item)
		}

		return nil
	})
}

func (s *ItemsService) GetById(ctx context.Context, itemId int) (models.Item, error) {
	var item models.Item
	return item, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		item, err = s.repo.GetById(ctx, itemId)
		if err != nil {
			return err
		}

		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
		if err != nil {
			return err
		}
		for i := range images {
			images[i].Filename = "/files/" + images[i].Filename
		}
		item.Images = images

		category, err := s.categoriesService.GetById(ctx, item.Category.Id)
		if err != nil {
			return err
		}
		item.Category = category

		return nil
	})
}

func (s *ItemsService) GetBySku(ctx context.Context, sku string) (models.Item, error) {
	var item models.Item
	return item, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		item, err = s.repo.GetBySku(ctx, sku)
		if err != nil {
			return err
		}

		colors, err := s.repo.GetColors(ctx, item.Id)
		if err != nil {
			return err
		}
		item.Colors = colors

		tags, err := s.repo.GetTags(ctx, item.Id)
		if err != nil {
			return err
		}
		item.Tags = tags

		images, err := s.repo.GetImages(ctx, item.Id)
		if err != nil {
			return err
		}
		for i := range images {
			images[i].Filename = "/files/" + images[i].Filename
		}
		item.Images = images

		category, err := s.categoriesService.GetById(ctx, item.Category.Id)
		if err != nil {
			return err
		}
		item.Category = category

		return nil
	})
}

func (s *ItemsService) GetByCategory(ctx context.Context, categoryId int) ([]models.Item, error) {
	var items []models.Item
	return items, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.categoriesService.Exist(ctx, categoryId)
		if !exist {
			return apperrors.ErrIdNotFound("category", categoryId)
		} else if err != nil {
			return err
		}

		ids, err := s.repo.GetByCategory(ctx, categoryId)
		if err != nil {
			return err
		}

		for _, id := range ids {
			item, err := s.repo.GetById(ctx, id)
			if err != nil {
				return err
			}
			colors, err := s.repo.GetColors(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Colors = colors

			tags, err := s.repo.GetTags(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Tags = tags

			images, err := s.repo.GetImages(ctx, item.Id)
			if err != nil {
				return err
			}
			for i := range images {
				images[i].Filename = "/files/" + images[i].Filename
			}
			item.Images = images

			category, err := s.categoriesService.GetById(ctx, item.Category.Id)
			if err != nil {
				return err
			}
			item.Category = category

			items = append(items, item)
		}

		return nil
	})
}

func (s *ItemsService) GetByTag(ctx context.Context, tag string) ([]models.Item, error) {
	var items []models.Item
	return items, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		ids, err := s.repo.GetByTag(ctx, tag)
		if err != nil {
			return err
		}

		for _, id := range ids {
			item, err := s.repo.GetById(ctx, id)
			if err != nil {
				return err
			}
			colors, err := s.repo.GetColors(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Colors = colors

			tags, err := s.repo.GetTags(ctx, item.Id)
			if err != nil {
				return err
			}
			item.Tags = tags

			images, err := s.repo.GetImages(ctx, item.Id)
			if err != nil {
				return err
			}
			for i := range images {
				images[i].Filename = "/files/" + images[i].Filename
			}
			item.Images = images

			category, err := s.categoriesService.GetById(ctx, item.Category.Id)
			if err != nil {
				return err
			}
			item.Category = category

			items = append(items, item)
		}

		return nil
	})
}

func (s *ItemsService) Update(ctx context.Context, item models.Item) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.Exist(ctx, item.Id)
		if !exist {
			return apperrors.ErrIdNotFound("item", item.Id)
		} else if err != nil {
			return err
		}

		// Check category exist
		exist, err = s.categoriesService.Exist(ctx, item.Category.Id)
		if !exist {
			return apperrors.ErrIdNotFound("category", item.Category.Id)
		} else if err != nil {
			return err
		}

		// Check colors exist
		for _, color := range item.Colors {
			exist, err := s.colorsService.Exist(ctx, color.Id)
			if !exist {
				return apperrors.ErrIdNotFound("color", item.Category.Id)
			} else if err != nil {
				return err
			}
		}

		// Check images exist
		for _, image := range item.Images {
			exist, err := s.imagesService.Exist(ctx, image.Id)
			if !exist {
				return apperrors.ErrIdNotFound("image", item.Category.Id)
			} else if err != nil {
				return err
			}
		}

		if err := s.repo.Update(ctx, item.Id, item.Name, item.Description, item.Category.Id, item.Price, item.Sku); err != nil {
			return err
		}

		// Update tags
		if err := s.repo.DeleteTags(ctx, item.Id); err != nil {
			return err
		}
		if len(item.Tags) > 0 {
			for _, tag := range item.Tags {
				if err := s.repo.LinkTag(ctx, item.Id, tag.Name); err != nil {
					return err
				}
			}
		}

		// Update colors
		if err := s.repo.DeleteColors(ctx, item.Id); err != nil {
			return err
		}
		for _, color := range item.Colors {
			if err := s.repo.LinkColor(ctx, item.Id, color.Id); err != nil {
				return err
			}
		}

		// Update images
		if err := s.repo.DeleteImages(ctx, item.Id); err != nil {
			return err
		}
		for _, image := range item.Images {
			if err := s.repo.LinkImage(ctx, item.Id, image.Id); err != nil {
				return err
			}
		}

		return nil
	})
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
