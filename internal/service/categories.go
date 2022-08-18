package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
)

type CategoriesService struct {
	repo          repository.Categories
	imagesService Images
}

func NewCategoriesService(repo repository.Categories, images Images) *CategoriesService {
	return &CategoriesService{
		repo:          repo,
		imagesService: images,
	}
}

func (s *CategoriesService) Create(ctx context.Context, name string, imageId int) (int, error) {
	var id int
	return id, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		category := models.Category{
			Name: name,
		}

		exist, err := s.imagesService.Exist(ctx, imageId)
		if !exist {
			return apperrors.ErrIdNotFound("image", imageId)
		} else if err != nil {
			return err
		}

		id, err = s.repo.Create(ctx, category)
		if err != nil {
			return err
		}

		if err := s.repo.LinkImage(ctx, id, imageId); err != nil {
			return err
		}

		return nil
	})
}

func (s *CategoriesService) Delete(ctx context.Context, categoryId int) error {
	return s.repo.Delete(ctx, categoryId)
}

func (s *CategoriesService) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	return categories, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		categories, err = s.repo.GetAll(ctx)
		if err != nil {
			return err
		}

		for i, category := range categories {
			image, err := s.repo.GetImage(ctx, category.Id)
			if err != nil {
				return err
			}
			image.Filename = "/files/" + image.Filename
			categories[i].Image = image
		}

		return nil
	})
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

		image, err := s.repo.GetImage(ctx, category.Id)
		if err != nil {
			return err
		}

		image.Filename = "/files/" + image.Filename
		category.Image = image

		return nil
	})
}

func (s *CategoriesService) Update(ctx context.Context, categoryId int, name string, imageId int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.repo.Exist(ctx, categoryId)
		if !exist {
			return apperrors.ErrIdNotFound("category", categoryId)
		} else if err != nil {
			return err
		}

		exist, err = s.imagesService.Exist(ctx, imageId)
		if !exist {
			return apperrors.ErrIdNotFound("image", imageId)
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
