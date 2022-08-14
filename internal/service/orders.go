package service

import (
	"context"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	apperrors "shop_backend/pkg/errors"
)

type OrdersService struct {
	repo repository.Orders

	usersService    Users
	deliveryService Delivery
	itemsService    Items
	colorsService   Colors
}

func NewOrdersService(repo repository.Orders, users Users, delivery Delivery, items Items, colors Colors) *OrdersService {
	return &OrdersService{
		repo:            repo,
		usersService:    users,
		deliveryService: delivery,
		itemsService:    items,
		colorsService:   colors,
	}
}

func (s *OrdersService) Create(ctx context.Context, order models.Order) (int, error) {
	var id int
	return id, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		exist, err := s.deliveryService.Exist(ctx, order.DeliveryId)
		if !exist {
			return apperrors.ErrIdNotFound("delivery", order.DeliveryId)
		} else if err != nil {
			return err
		}

		id, err = s.repo.Create(ctx, order.UserId, order.DeliveryId)
		if err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := s.repo.LinkItem(ctx, id, item.Id, item.ColorId, item.Quantity); err != nil {
				return err
			}
		}

		return nil
	})
}
