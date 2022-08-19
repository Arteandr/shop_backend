package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
	"strconv"
	"strings"
)

func (h *Handler) InitDeliveryRoutes(api *gin.RouterGroup) {
	delivery := api.Group("/delivery")
	{
		admins := delivery.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.completedIdentify, h.createDelivery)
			admins.GET("/:id", h.completedIdentify, h.getDeliveryById)
			admins.PUT("/:id", h.completedIdentify, h.updateDelivery)
			admins.DELETE("/:id", h.completedIdentify, h.deleteDelivery)
		}
		delivery.GET("/all", h.getAllDelivery)
	}
}

type createDeliveryInput struct {
	Name        string  `json:"name" binding:"required"`
	CompanyName string  `json:"companyName" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

func (i *createDeliveryInput) isValid() error {
	if len(i.Name) < 1 || len(i.Name) > 30 {
		return errors.New("wrong name length")
	}

	if len(i.CompanyName) < 1 || len(i.CompanyName) > 30 {
		return errors.New("wrong company name length")
	}
	i.CompanyName = strings.ToLower(strings.TrimSpace(i.CompanyName))

	return nil
}

// @Summary Create a new delivery
// @Security UsersAuth
// @Security AdminAuth
// @Tags delivery-actions
// @Description create a new delivery
// @Accept json
// @Produce json
// @Param input body createDeliveryInput true "delivery info"
// @Success 200 {object} IdResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delivery/create [post]
func (h *Handler) createDelivery(ctx *gin.Context) {
	var body createDeliveryInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	delivery := models.Delivery{
		Name:        body.Name,
		CompanyName: body.CompanyName,
		Price:       body.Price,
	}

	id, err := h.services.Delivery.Create(ctx.Request.Context(), delivery)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, IdResponse{Id: id})
}

// @Summary Get delivery by id
// @Security UsersAuth
// @Security AdminAuth
// @Tags delivery-actions
// @Description get delivery by id
// @Accept json
// @Produce json
// @Param id path int true "delivery id"
// @Success 200 {object} models.Delivery
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delivery/{id} [get]
func (h *Handler) getDeliveryById(ctx *gin.Context) {
	strDeliveryId := ctx.Param("id")
	deliveryId, err := strconv.Atoi(strDeliveryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	delivery, err := h.services.Delivery.GetById(ctx.Request.Context(), deliveryId)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)
			return
		}
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// @Summary Get all delivery
// @Tags delivery-actions
// @Description get all delivery
// @Accept json
// @Produce json
// @Success 200 {array} models.Delivery
// @Failure 500 {object} ErrorResponse
// @Router /delivery/all [get]
func (h *Handler) getAllDelivery(ctx *gin.Context) {
	delivery, err := h.services.Delivery.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// @Summary Update delivery by id
// @Security UsersAuth
// @Security AdminAuth
// @Tags delivery-actions
// @Description update delivery by id
// @Accept json
// @Produce json
// @Param id path int true "delivery id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delivery/{id} [put]
func (h *Handler) updateDelivery(ctx *gin.Context) {
	var body createDeliveryInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	strDeliveryId := ctx.Param("id")
	deliveryId, err := strconv.Atoi(strDeliveryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	delivery := models.Delivery{
		Id:          deliveryId,
		Name:        body.Name,
		CompanyName: body.CompanyName,
		Price:       body.Price,
	}

	if err := h.services.Delivery.Update(ctx.Request.Context(), delivery); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Delete delivery by id
// @Security UsersAuth
// @Security AdminAuth
// @Tags delivery-actions
// @Description delete delivery by id
// @Accept json
// @Produce json
// @Param id path int true "delivery id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delivery/{id} [delete]
func (h *Handler) deleteDelivery(ctx *gin.Context) {
	strDeliveryId := ctx.Param("id")
	deliveryId, err := strconv.Atoi(strDeliveryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	if err := h.services.Delivery.Delete(ctx, deliveryId); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
