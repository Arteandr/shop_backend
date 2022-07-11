package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
)

func (h *Handler) InitColorsRoutes(api *gin.RouterGroup) {
	colors := api.Group("/colors")
	{
		colors.POST("/create", h.createColor)
	}
}

func (h *Handler) createColor(ctx *gin.Context) {
	var color models.Color
	if err := ctx.BindJSON(&color); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	colorId, err := h.services.Colors.Create(color.Name, color.Hex, color.Price)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"colorId": colorId,
	})
}
