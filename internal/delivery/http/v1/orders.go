package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
	"strconv"
)

func (h *Handler) InitOrdersRoutes(api *gin.RouterGroup) {
	orders := api.Group("/orders", h.userIdentity)
	{
		admin := orders.Group("/", h.adminIdentify)
		{
			admin.DELETE("/:id", h.deleteOrder)
			admin.PUT("/:id", h.updateOrderStatus)
			admin.GET("/statuses/all", h.getAllOrderStatuses)
		}
		orders.GET("/me/all", h.completedIdentify, h.getAllOrders)
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

// @Summary Delete order
// @Security UsersAuth
// @Security AdminAuth
// @Tags orders-actions
// @Description delete order by id
// @Accept json
// @Produce json
// @Param id path int true "order id"
// @Success 200 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [delete]
func (h *Handler) deleteOrder(ctx *gin.Context) {
	strOrderId := ctx.Param("id")
	orderId, err := strconv.Atoi(strOrderId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	if err := h.services.Orders.Delete(ctx.Request.Context(), orderId); err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)
			return
		}
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get all orders
// @Security UsersAuth
// @Tags orders-actions
// @Description get all orders
// @Accept json
// @Produce json
// @Success 200 {array} models.ServiceOrder
// @Failure 500 {object} ErrorResponse
// @Router /orders/me/all [get]
func (h *Handler) getAllOrders(ctx *gin.Context) {
	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	orders, err := h.services.Orders.GetAllByUserId(ctx, userId)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// @Summary Get all order statuses
// @Security UsersAuth
// @Security AdminAuth
// @Tags orders-actions
// @Description get all order statuses
// @Accept json
// @Produce json
// @Success 200 {array} models.OrderStatus
// @Failure 500 {object} ErrorResponse
// @Router /orders/statuses/all [get]
func (h *Handler) getAllOrderStatuses(ctx *gin.Context) {
	statuses, err := h.services.Orders.GetAllStatuses(ctx.Request.Context())
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, statuses)
}

type updateOrderStatusInput struct {
	StatusId int `json:"statusId" binding:"required"`
}

// @Summary Update order status
// @Security UsersAuth
// @Security AdminAuth
// @Tags orders-actions
// @Description update order status by id
// @Accept json
// @Produce json
// @Success 200 ""
// @Param id path int true "order id"
// @Param input body updateOrderStatusInput true "status info"
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [put]
func (h *Handler) updateOrderStatus(ctx *gin.Context) {
	strOrderId := ctx.Param("id")
	orderId, err := strconv.Atoi(strOrderId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	var body updateOrderStatusInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	if err := h.services.Orders.UpdateStatus(ctx.Request.Context(), orderId, body.StatusId); err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)
			return
		}
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
