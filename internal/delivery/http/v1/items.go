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
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	CategoryId  int          `json:"categoryId" binding:"required"`
	Tags        []string     `json:"tags" binding:"required"`
	Colors      models.Color `json:"colors" binding:"required"`
}

func (h *Handler) createItem(ctx *gin.Context) {
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.services.Items.Create(body.Name, body.Description, body.CategoryId, body.Tags, time.Now())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
