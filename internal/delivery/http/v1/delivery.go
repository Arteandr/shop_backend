package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"strconv"
	"strings"
)

func (h *Handler) InitDeliveryRoutes(api *gin.RouterGroup) {
	delivery := api.Group("/delivery")
	{
		admins := delivery.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.createDelivery)
			admins.GET("/:id", h.getDeliveryById)
		}
		delivery.PUT("/:id", h.updateDelivery)
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
// @Success 200 {object} CreatDeliveryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /delivery/create [post]
func (h *Handler) createDelivery(ctx *gin.Context) {
	var body createDeliveryInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	delivery := models.Delivery{
		Name:        body.Name,
		CompanyName: body.CompanyName,
		Price:       body.Price,
	}

	id, err := h.services.Delivery.Create(ctx.Request.Context(), delivery)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, CreatDeliveryResponse{
		Id: id,
	})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	delivery, err := h.services.Delivery.GetById(ctx.Request.Context(), deliveryId)
	if err != nil && errors.Is(err, models.ErrDeliveryNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

func (h *Handler) updateDelivery(ctx *gin.Context) {
	var body createDeliveryInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	strDeliveryId := ctx.Param("id")
	deliveryId, err := strconv.Atoi(strDeliveryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	delivery := models.Delivery{
		Id:          deliveryId,
		Name:        body.Name,
		CompanyName: body.CompanyName,
		Price:       body.Price,
	}

	if err := h.services.Delivery.Update(ctx.Request.Context(), delivery); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
