package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"strconv"
)

func (h *Handler) InitColorsRoutes(api *gin.RouterGroup) {
	colors := api.Group("/colors")
	{
		colors.POST("/create", h.createColor)
		colors.DELETE("/:id", h.deleteColor)
	}
}

type CreateColorResult struct {
	ColorId int `json:"colorId"`
}

// @Summary Create a new color
// @Tags colors-actions
// @Description create a new color
// @Accept json
// @Produce json
// @Param input body models.Color true "input body"
// @Success 200 {object} CreateColorResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/create [post]
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

	ctx.JSON(http.StatusOK, CreateColorResult{ColorId: colorId})
}

// @Summary Delete colors
// @Tags colors-actions
// @Description delete color by id
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/{id} [delete]
func (h *Handler) deleteColor(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Colors.Delete(colorId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
