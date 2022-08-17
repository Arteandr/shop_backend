package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
)

func (h *Handler) InitOrdersRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders", h.userIdentity)
	{
		orders.POST("/create", h.completedIdentify, h.createOrder)
	}
}

type createOrderInput struct {
	Items      []models.OrderItem `json:"items" binding:"required"`
	DeliveryId int                `json:"deliveryId" binding:"required"`
}

func (i *createOrderInput) isValid() error {
	if len(i.Items) < 1 {
		return errors.New("wrong items length")
	}
	for _, item := range i.Items {
		if item.Id < 1 {
			return fmt.Errorf("wrong item id %d", item.Id)
		}
		if item.Quantity < 1 {
			return fmt.Errorf("wrong quantity %d", item.Quantity)
		}
		if item.ColorId < 1 {
			return fmt.Errorf("wrong color id %d", item.ColorId)
		}
	}
	if i.DeliveryId < 1 {
		return fmt.Errorf("wrong delivery id %d", i.DeliveryId)
	}

	return nil
}

// @Summary Create a new order
// @Security UsersAuth
// @Tags orders-actions
// @Description create a new order
// @Accept json
// @Produce json
// @Param input body createOrderInput true "input body"
// @Success 201 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/create [post]
func (h *Handler) createOrder(ctx *gin.Context) {
	var body createOrderInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	order := models.Order{
		UserId:     userId,
		DeliveryId: body.DeliveryId,
		Items:      body.Items,
	}

	id, err := h.services.Orders.Create(ctx, order)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)
			return
		}
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, IdResponse{Id: id})
}
