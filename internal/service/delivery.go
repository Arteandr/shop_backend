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
	var id int
	return id, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		companyExist, err := s.repo.ExistCompany(ctx, delivery.CompanyName)
		if err != nil {
			return err
		}

		if !companyExist {
			if err := s.repo.CreateCompany(ctx, delivery.CompanyName); err != nil {
				return err
			}
		}

		id, err = s.repo.Create(ctx, delivery)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *DeliveryService) Delete(ctx context.Context, deliveryId int) error {
	return s.repo.Delete(ctx, deliveryId)
}

func (s *DeliveryService) Update(ctx context.Context, delivery models.Delivery) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		companyExist, err := s.repo.ExistCompany(ctx, delivery.CompanyName)
		if err != nil {
			return err
		}

		if !companyExist {
			if err := s.repo.CreateCompany(ctx, delivery.CompanyName); err != nil {
				return err
			}
		}

		if err := s.repo.Update(ctx, delivery); err != nil {
			return err
		}

		return nil
	})
}

func (s *DeliveryService) GetById(ctx context.Context, deliveryId int) (models.Delivery, error) {
	delivery, err := s.repo.GetById(ctx, deliveryId)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}

func (s *DeliveryService) GetAll(ctx context.Context) ([]models.Delivery, error) {
	return s.repo.GetAll(ctx)
}

func (s *DeliveryService) Exist(ctx context.Context, deliveryId int) (bool, error) {
	return s.repo.Exist(ctx, deliveryId)
}
