package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
)

func (h *Handler) InitCategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		categories.POST("/create", h.createCategory)
	}
}

type CreateCategoryResult struct {
	CategoryId int `json:"categoryId"`
}

// @Summary Create a new category
// @Tags categories-actions
// @Description create a new category
// @Accept json
// @Produce json
// @Param input body models.Category true "input body"
// @Success 200 {object} CreateCategoryResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/create [post]
func (h *Handler) createCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.BindJSON(&category); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	categoryId, err := h.services.Categories.Create(category.Name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, CreateCategoryResult{CategoryId: categoryId})
}
