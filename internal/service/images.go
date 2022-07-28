package service

import (
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type ImagesService struct {
	repo repository.Images
}

func NewImagesService(repo repository.Images) *ImagesService {
	return &ImagesService{repo: repo}
}

func (s *ImagesService) Upload(filename string) (int, error) {
	return s.repo.Upload(filename)
}

func (s *ImagesService) GetAll() ([]models.Image, error) {
	return s.repo.GetAll()
}
