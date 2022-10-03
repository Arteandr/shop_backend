package service

import (
	"context"
	"errors"
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
	mailsService    Mails
}

func NewOrdersService(repo repository.Orders, users Users, delivery Delivery, items Items, colors Colors, mails Mails) *OrdersService {
	return &OrdersService{
		repo:            repo,
		usersService:    users,
		deliveryService: delivery,
		itemsService:    items,
		colorsService:   colors,
		mailsService:    mails,
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

		id, err = s.repo.Create(ctx, order.UserId, order.DeliveryId, order.Comment)
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

		for i := range orders {
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
				Comment:   orders[i].Comment,
				CreatedAt: orders[i].CreatedAt,
			}

			fullOrders = append(fullOrders, o)
		}

		return nil
	})
}

func (s *OrdersService) GetAll(ctx context.Context) ([]models.ServiceOrder, error) {
	var fullOrders []models.ServiceOrder
	return fullOrders, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		orders, err := s.repo.GetAll(ctx)
		if err != nil {
			return err
		}

		for i := range orders {
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
				Comment:   orders[i].Comment,
				CreatedAt: orders[i].CreatedAt,
			}

			fullOrders = append(fullOrders, o)
		}

		return nil
	})
}
func (s *OrdersService) GetById(ctx context.Context, orderId int) (models.ServiceOrder, error) {
	var order models.ServiceOrder
	return order, s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		o, err := s.repo.GetById(ctx, orderId)
		if err != nil {
			return err
		}

		items, err := s.repo.GetItems(ctx, o.Id)
		if err != nil {
			return err
		}

		delivery, err := s.deliveryService.GetById(ctx, o.DeliveryId)
		if err != nil {
			return err
		}

		status, err := s.repo.GetStatus(ctx, o.StatusId)
		if err != nil {
			return err
		}

		order = models.ServiceOrder{
			Id:        o.Id,
			Status:    status,
			UserId:    o.UserId,
			Items:     items,
			Delivery:  delivery,
			Comment:   o.Comment,
			CreatedAt: o.CreatedAt,
		}

		return nil
	})
}

func (s *OrdersService) GetAllStatuses(ctx context.Context) ([]models.OrderStatus, error) {
	return s.repo.GetAllStatuses(ctx)
}

func (s *OrdersService) GetAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	methods, err := s.repo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	for i := range methods {
		if methods[i].Logo != nil {
			if len(*methods[i].Logo) > 0 {
				*methods[i].Logo = "/files/" + *methods[i].Logo
			}
		}
	}

	return methods, nil
}

func (s *OrdersService) UpdateStatus(ctx context.Context, orderId, statusId int) error {
	return s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		var exist bool
		var err error

		exist, err = s.repo.ExistStatus(ctx, statusId)
		if !exist {
			return apperrors.ErrIdNotFound("status", statusId)
		} else if err != nil {
			return err
		}

		status, err := s.repo.GetStatus(ctx, statusId)
		if err != nil {
			return err
		}

		exist, err = s.repo.Exist(ctx, orderId)
		if !exist {
			return apperrors.ErrIdNotFound("order", orderId)
		} else if err != nil {
			return err
		}

		order, err := s.repo.GetById(ctx, orderId)
		if err != nil {
			return err
		}

		user, err := s.usersService.GetMe(ctx, order.UserId)
		if err != nil {
			return err
		}

		if len(*user.FirstName) <= 0 || len(*user.LastName) <= 0 {
			return errors.New("user don't have firstName or lastName")
		}

		if err := s.repo.UpdateStatus(ctx, orderId, statusId); err != nil {
			return err
		}

		if err := s.mailsService.UpdateStatus(ctx, user.Login, user.Email, *user.FirstName, *user.LastName, status); err != nil {
			return err
		}

		return nil
	})
}
