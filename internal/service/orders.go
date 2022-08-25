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

func (s *OrdersService) Delete(ctx context.Context, orderId int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := s.Exist(ctx, orderId)
		if !exist {
			return apperrors.ErrIdNotFound("order", orderId)
		} else if err != nil {
			return err
		}

		if err := s.repo.Delete(ctx, orderId); err != nil {
			return err
		}

		return nil
	})
}

func (s *OrdersService) Exist(ctx context.Context, orderId int) (bool, error) {
	return s.repo.Exist(ctx, orderId)
}

func (s *OrdersService) GetAllByUserId(ctx context.Context, userId int) ([]models.ServiceOrder, error) {
	var fullOrders []models.ServiceOrder
	return fullOrders, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		orders, err := s.repo.GetAllByUserId(ctx, userId)
		if err != nil {
			return err
		}

		for i, _ := range orders {
			items, err := s.repo.GetItems(ctx, orders[i].Id)
			if err != nil {
				return err
			}

			delivery, err := s.deliveryService.GetById(ctx, orders[i].DeliveryId)
			if err != nil {
				return err
			}

			status, err := s.repo.GetStatus(ctx, orders[i].StatusId)
			if err != nil {
				return err
			}

			o := models.ServiceOrder{
				Id:        orders[i].Id,
				Status:    status,
				UserId:    orders[i].UserId,
				Items:     items,
				Delivery:  delivery,
				CreatedAt: orders[i].CreatedAt,
			}

			fullOrders = append(fullOrders, o)
		}

		return nil
	})
}
