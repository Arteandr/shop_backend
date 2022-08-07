package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"strconv"
)

func (h *Handler) InitCategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		admins := categories.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.createCategory)
			admins.DELETE("/:id", h.deleteCategory)
		}
		categories.GET("/", h.getAllCategories)
		categories.GET("/:id", h.getCategoryById)
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

// @Summary Delete category
// @Tags categories-actions
// @Description delete category by id
// @Accept json
// @Produce json
// @Param id path int true "category id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *Handler) deleteCategory(ctx *gin.Context) {
	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Categories.Delete(categoryId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get category by id
// @Tags categories-actions
// @Description get category by id
// @Accept json
// @Produce json
// @Success 200 {object} models.Category
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *Handler) getCategoryById(ctx *gin.Context) {
	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if exist, err := h.services.Categories.Exist(categoryId); !exist {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("wrong category %d id", categoryId)})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.services.Categories.GetById(categoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// @Summary Get all categories
// @Tags categories-actions
// @Description get all categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} ErrorResponse
// @Router /categories/ [get]
func (h *Handler) getAllCategories(ctx *gin.Context) {
	categories, err := h.services.Categories.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
