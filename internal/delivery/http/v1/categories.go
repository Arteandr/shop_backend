package v1

import (
	"errors"
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
			admins.PUT("/:id", h.updateCategory)
		}
		categories.GET("/", h.getAllCategories)
		categories.GET("/:id", h.getCategoryById)
	}
}

// @Summary Create a new category
// @Security UsersAuth
// @Security AdminAuth
// @Tags categories-actions
// @Description create a new category
// @Accept json
// @Produce json
// @Param input body models.Category true "input body"
// @Success 200 {object} IdResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/create [post]
func (h *Handler) createCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.BindJSON(&category); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	categoryId, err := h.services.Categories.Create(ctx.Request.Context(), category.Name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, IdResponse{Id: categoryId})
}

// @Summary Delete category
// @Security UsersAuth
// @Security AdminAuth
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

	if err := h.services.Categories.Delete(ctx.Request.Context(), categoryId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

type updateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (u *updateCategoryInput) isValid() error {
	if len(u.Name) < 1 || len(u.Name) > 15 {
		return errors.New("wrong name length")
	}

	return nil
}

// @Summary Update category
// @Security UsersAuth
// @Security AdminAuth
// @Tags categories-actions
// @Description update category by id
// @Accept json
// @Produce json
// @Param id path int true "category id"
// @Param input body updateCategoryInput true "name info"
// @Success 200 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *Handler) updateCategory(ctx *gin.Context) {
	var body updateCategoryInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValid(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if exist, err := h.services.Categories.Exist(ctx.Request.Context(), categoryId); err != nil || !exist {
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("wrong category id %d", categoryId)})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Categories.Update(ctx.Request.Context(), categoryId, body.Name); err != nil {
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

	if exist, err := h.services.Categories.Exist(ctx.Request.Context(), categoryId); !exist {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("wrong category %d id", categoryId)})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.services.Categories.GetById(ctx.Request.Context(), categoryId)
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
	categories, err := h.services.Categories.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
