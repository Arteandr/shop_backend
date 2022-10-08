package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "shop_backend/pkg/errors"
)

func (h *Handler) InitCategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		admins := categories.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.completedIdentify, h.createCategory)
			admins.DELETE("/:id", h.completedIdentify, h.deleteCategory)
			admins.PUT("/:id", h.completedIdentify, h.updateCategory)
		}
		categories.GET("/", h.getAllCategories)
		categories.GET("/:id", h.getCategoryById)
	}
}

type createCategoryInput struct {
	Name    string `json:"name" binding:"required"`
	ImageId int    `json:"imageId" binding:"required"`
}

func (i *createCategoryInput) isValid() error {
	if len(i.Name) < 1 || len(i.Name) > 30 {
		return errors.New("wrong category name")
	}

	if i.ImageId < 1 {
		return errors.New("wrong image id")
	}

	return nil
}

// @Summary Create a new category
// @Security UsersAuth
// @Security AdminAuth
// @Tags categories-actions
// @Description create a new category
// @Accept json
// @Produce json
// @Param input body createCategoryInput true "input body"
// @Success 201 {object} IdResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/create [post]
func (h *Handler) createCategory(ctx *gin.Context) {
	var body createCategoryInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	categoryId, err := h.services.Categories.Create(ctx.Request.Context(), body.Name, body.ImageId)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)
			return
		}

		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, IdResponse{Id: categoryId})
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
// @Failure 400,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *Handler) deleteCategory(ctx *gin.Context) {
	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	if err := h.services.Categories.Delete(ctx.Request.Context(), categoryId); err != nil {
		if errors.Is(err, apperrors.ErrViolatesKey) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.Status(http.StatusOK)
}

type updateCategoryInput struct {
	Name    string `json:"name" binding:"required"`
	ImageId int    `json:"imageId" binding:"required"`
}

func (i *updateCategoryInput) isValid() error {
	if len(i.Name) < 1 || len(i.Name) > 15 {
		return errors.New("wrong name length")
	}

	if i.ImageId < 1 {
		return errors.New("wrong image id")
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
	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	var body updateCategoryInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)

		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	if err := h.services.Categories.Update(ctx.Request.Context(), categoryId, body.Name, body.ImageId); err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

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
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	category, err := h.services.Categories.GetById(ctx.Request.Context(), categoryId)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

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
// @Router /categories [get]
func (h *Handler) getAllCategories(ctx *gin.Context) {
	categories, err := h.services.Categories.GetAll(ctx.Request.Context())
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, categories)
}
