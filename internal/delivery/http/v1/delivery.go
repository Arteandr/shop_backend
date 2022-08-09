package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"strings"
)

func (h *Handler) InitDeliveryRoutes(api *gin.RouterGroup) {
	delivery := api.Group("/delivery")
	{
		//admins := delivery.Group("/", h.userIdentity, h.adminIdentify)
		//{
		//	admins.POST("/", h.createDelivery)
		//}
		delivery.POST("/", h.createDelivery)

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
