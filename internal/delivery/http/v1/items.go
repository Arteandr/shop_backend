package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"time"
)

func (h *Handler) InitItemsRoutes(api *gin.RouterGroup) {
	items := api.Group("/items")
	{
		items.POST("/create", h.createItem)
	}
}

type createItemInput struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description" binding:"required"`
	CategoryId  int            `json:"categoryId" binding:"required"`
	Tags        []string       `json:"tags" binding:"required"`
	Colors      []models.Color `json:"colors" binding:"required"`
}

func (h *Handler) createItem(ctx *gin.Context) {
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	itemId, err := h.services.Items.Create(body.Name, body.Description, body.CategoryId, body.Tags, time.Now())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := 0; i < len(body.Colors); i++ {
		color := body.Colors[i]
		colorId, err := h.services.Colors.Create(color.Name, color.Hex, color.Price)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		if err := h.services.Items.LinkColor(itemId, colorId); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, models.Item{
		Id:          itemId,
		Name:        body.Name,
		Description: body.Description,
		CategoryId:  body.CategoryId,
		Tags:        body.Tags,
		CreatedAt:   time.Now(),
		Colors:      body.Colors,
	})
}
