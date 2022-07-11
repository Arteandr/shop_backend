package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
)

func (h *Handler) InitColorsRoutes(api *gin.RouterGroup) {
	colors := api.Group("/colors")
	{
		colors.POST("/create", h.createColor)
		colors.POST("/exist", h.checkExist)
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

func (h *Handler) checkExist(ctx *gin.Context) {
	var b struct {
		Id int `json:"colorId"`
	}
	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	exist, err := h.services.Colors.Exist(b.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	fmt.Println("Exist:", exist)

	ctx.JSON(http.StatusOK, exist)
}
