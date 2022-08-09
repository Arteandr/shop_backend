package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
)

type DeliveryService struct {
	repo repository.Delivery
}

func NewDeliveryService(repo repository.Delivery) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) Create(ctx context.Context, delivery models.Delivery) (int, error) {
	companyExist, err := s.repo.ExistCompany(ctx, delivery.CompanyName)
	if err != nil {
		return 0, err
	}

	if !companyExist {
		if err := s.repo.CreateCompany(ctx, delivery.CompanyName); err != nil {
			return 0, err
		}
	}

	id, err := s.repo.Create(ctx, delivery)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *DeliveryService) GetById(ctx context.Context, deliveryId int) (models.Delivery, error) {
	delivery, err := s.repo.GetById(ctx, deliveryId)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}
