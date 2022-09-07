package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
)

func (h *Handler) InitItemsRoutes(api *gin.RouterGroup) {
	items := api.Group("/items")
	{
		admins := items.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.completedIdentify, h.createItem)
			admins.PUT("/:id", h.completedIdentify, h.updateItems)
			admins.DELETE("/", h.completedIdentify, h.deleteItems)
			admins.GET("/all", h.completedIdentify, h.sort("created_at", ASC), h.getAllItems)
		}

		items.GET("/new", h.getNewItems)
		items.GET("/:id", h.getItemById)
		items.GET("/sku/:sku", h.getItemBySku)
		items.GET("/category/:id", h.getItemsByCategory)
		items.GET("/tag/:name", h.getItemsByTag)
	}
}

type createItemInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CategoryId  int      `json:"categoryId" binding:"required"`
	Tags        []string `json:"tags"`
	ColorsId    []int    `json:"colors" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Sku         string   `json:"sku" binding:"required"`
	ImagesId    []int    `json:"images" binding:"required"`
}

func (i *createItemInput) isValid() error {
	if len(i.Name) < 1 || len(i.Name) > 50 {
		return errors.New("wrong name length")
	}

	if len(i.Description) < 1 || len(i.Description) > 255 {
		return errors.New("wrong description length")
	}

	if i.CategoryId < 1 {
		return errors.New("wrong category id")
	}

	if len(i.ColorsId) < 1 {
		return errors.New("wrong colors id")
	}

	if i.Price < 0 {
		return errors.New("wrong price")
	}

	if len(i.Sku) < 1 || len(i.Sku) > 20 {
		return errors.New("wrong sku")
	}

	if len(i.ImagesId) < 1 {
		return errors.New("the item must contain at least 1 image")
	}

	return nil
}

// @Summary Create a new item
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description create a new item
// @Accept json
// @Produce json
// @Param input body createItemInput true "input body"
// @Success 201 {object} models.Item
// @Failure 400,404,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/create [post]
func (h *Handler) createItem(ctx *gin.Context) {
	// Body binding
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)

		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	// Create item
	c := make([]models.Color, len(body.ColorsId))
	for i, colorId := range body.ColorsId {
		c[i] = models.Color{Id: colorId}
	}

	t := make([]models.Tag, len(body.Tags))
	for i, tag := range body.Tags {
		t[i] = models.Tag{Name: tag}
	}

	imgs := make([]models.Image, len(body.ImagesId))
	for i, imgId := range body.ImagesId {
		imgs[i] = models.Image{Id: imgId}
	}

	i := models.Item{
		Name:        body.Name,
		Description: body.Description,
		Category:    models.Category{Id: body.CategoryId},
		Sku:         body.Sku,
		Price:       body.Price,
		Colors:      c,
		Tags:        t,
		Images:      imgs,
	}
	item, err := h.services.Items.Create(ctx.Request.Context(), i)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// @Summary Get all items
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description get all items with sort options
// @Accept json
// @Produce json
// @Param sort_by query string false "sort field" Enums(id,name,description,category_id,price,sku,created_at)
// @Param sort_order query string false "sort order" Enums(asc,desc)
// @Success 200 {array} models.Item
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/all [get]
func (h *Handler) getAllItems(ctx *gin.Context) {
	sortOptions, err := getSortOptions(ctx)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrSortOptions)

		return
	}

	items, err := h.services.Items.GetAll(ctx.Request.Context(), sortOptions)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, items)
}

// @Summary Get new items
// @Tags items-actions
// @Description get new items
// @Accept json
// @Produce json
// @Success 200 {array} models.Item
// @Failure 500 {object} ErrorResponse
// @Router /items/new [get]
func (h *Handler) getNewItems(ctx *gin.Context) {
	items, err := h.services.Items.GetNew(ctx.Request.Context())
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, items)
}

// @Summary Get item by ID
// @Tags items-actions
// @Description get item by id
// @Accept json
// @Produce json
// @Param id path int true "item id"
// @Success 200 {object} models.Item
// @Failure 400,404,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/{id} [get]
func (h *Handler) getItemById(ctx *gin.Context) {
	strItemId := ctx.Param("id")
	itemId, err := strconv.Atoi(strItemId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	item, err := h.services.Items.GetById(ctx.Request.Context(), itemId)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, item)
}

// @Summary Get items with category
// @Tags items-actions
// @Description get all items with provided category id
// @Accept json
// @Produce json
// @Param id path int true "category id"
// @Success 200 {array} models.Item
// @Failure 400,404,409 {object} ErrorResponse
// @Router /items/category/{id} [get]
func (h *Handler) getItemsByCategory(ctx *gin.Context) {
	strCategoryId := ctx.Param("id")
	categoryId, err := strconv.Atoi(strCategoryId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	items, err := h.services.Items.GetByCategory(ctx.Request.Context(), categoryId)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, items)
}

// @Summary Get items with tag
// @Tags items-actions
// @Description get all items with provided tag id
// @Accept json
// @Produce json
// @Param id path int true "tag id"
// @Success 200 {array} models.Item
// @Failure 400,404,409 {object} ErrorResponse
// @Router /items/tag/{id} [get]
func (h *Handler) getItemsByTag(ctx *gin.Context) {
	tag := ctx.Param("name")
	if len(tag) < 1 {
		err := errors.New("tag must contain at least 1 character")
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	items, err := h.services.Items.GetByTag(ctx.Request.Context(), tag)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, items)
}

// @Summary Get item by SKU
// @Tags items-actions
// @Description get item by sku
// @Accept json
// @Produce json
// @Param sku path string true "item sku"
// @Success 200 {object} models.Item
// @Failure 400,404,409 {object} ErrorResponse
// @Router /items/sku/{sku} [get]
func (h *Handler) getItemBySku(ctx *gin.Context) {
	sku := ctx.Param("sku")
	if len(sku) < 1 {
		err := errors.New("sku must contain at least 1 character")
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	item, err := h.services.Items.GetBySku(ctx.Request.Context(), sku)
	if err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, item)
}

type deleteItemsInput struct {
	ItemsId []int `json:"itemsId" binding:"required"`
}

func (i *deleteItemsInput) isValid() error {
	if len(i.ItemsId) < 1 {
		return fmt.Errorf("wrong items id length %d", len(i.ItemsId))
	}

	return nil
}

// @Summary Delete items
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description delete items by id
// @Accept json
// @Produce json
// @Param input body deleteItemsInput true "items id info"
// @Success 200 ""
// @Failure 400,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items [delete]
func (h *Handler) deleteItems(ctx *gin.Context) {
	var body deleteItemsInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)

		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	if err := h.services.Items.Delete(ctx.Request.Context(), body.ItemsId); err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.Status(http.StatusOK)
}

type updateItemInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CategoryId  int      `json:"categoryId" binding:"required"`
	Tags        []string `json:"tags"`
	ColorsId    []int    `json:"colors" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Sku         string   `json:"sku" binding:"required"`
	ImagesId    []int    `json:"images" binding:"required"`
}

func (i *updateItemInput) isValid() error {
	if len(i.Name) < 1 || len(i.Name) > 50 {
		return errors.New("wrong name length")
	}

	if len(i.Description) < 1 || len(i.Description) > 255 {
		return errors.New("wrong description length")
	}

	if i.CategoryId < 1 {
		return errors.New("wrong category id")
	}

	if len(i.ColorsId) < 1 {
		return errors.New("wrong colors id")
	}

	if i.Price < 0 {
		return errors.New("wrong price")
	}

	if len(i.Sku) < 1 || len(i.Sku) > 20 {
		return errors.New("wrong sku")
	}

	if len(i.ImagesId) < 1 {
		return errors.New("the item must contain at least 1 image")
	}

	return nil
}

// @Summary Update item
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description update item
// @Accept json
// @Produce json
// @Param id path string true "item id"
// @Param input body updateItemInput true "item body"
// @Success 200 ""
// @Failure 400,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/{id} [put]
func (h *Handler) updateItems(ctx *gin.Context) {
	strItemId := ctx.Param("id")
	itemId, err := strconv.Atoi(strItemId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)

		return
	}

	var body updateItemInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)

		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)

		return
	}

	// Update item
	c := make([]models.Color, len(body.ColorsId))
	for i, colorId := range body.ColorsId {
		c[i] = models.Color{Id: colorId}
	}

	t := make([]models.Tag, len(body.Tags))
	for i, tag := range body.Tags {
		t[i] = models.Tag{Name: tag}
	}

	imgs := make([]models.Image, len(body.ImagesId))
	for i, imgId := range body.ImagesId {
		imgs[i] = models.Image{Id: imgId}
	}

	i := models.Item{
		Id:          itemId,
		Name:        body.Name,
		Description: body.Description,
		Category:    models.Category{Id: body.CategoryId},
		Sku:         body.Sku,
		Price:       body.Price,
		Colors:      c,
		Tags:        t,
		Images:      imgs,
	}
	if err := h.services.Items.Update(ctx.Request.Context(), i); err != nil {
		if errors.As(err, &apperrors.IdNotFound{}) {
			NewError(ctx, http.StatusNotFound, err)

			return
		} else if errors.As(err, &apperrors.UniqueValue{}) {
			NewError(ctx, http.StatusConflict, err)

			return
		}

		NewError(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.Status(http.StatusOK)
}
