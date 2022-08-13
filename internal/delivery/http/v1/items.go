package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
	"strconv"
)

func (h *Handler) InitItemsRoutes(api *gin.RouterGroup) {
	items := api.Group("/items")
	{
		admins := items.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/create", h.createItem)
			admins.PUT("/:id", h.updateItems)
			admins.DELETE("/", h.deleteItems)
			admins.GET("/all", h.sort("created_at", ASC), h.getAllItems)

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

// @Summary Create a new item
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description create a new item
// @Accept json
// @Produce json
// @Param input body createItemInput true "input body"
// @Success 200 {object} models.Item
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/create [post]
func (h *Handler) createItem(ctx *gin.Context) {
	// Body binding
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	//// Check if exist category
	//if exist, err := h.services.Categories.Exist(ctx.Request.Context(), body.CategoryId); err != nil || !exist {
	//	NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("category", body.CategoryId))
	//	return
	//}

	//// Check if colors is existed
	//for _, colorId := range body.ColorsId {
	//	if exist, err := h.services.Colors.Exist(ctx.Request.Context(), colorId); err != nil || !exist {
	//		NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("color", colorId))
	//		return
	//	}
	//}

	// Check if at least one image
	if len(body.ImagesId) < 1 {
		err := errors.New("the item must contain at least 1 image")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Check if images is existed
	for _, imageId := range body.ImagesId {
		if exist, err := h.services.Images.Exist(ctx.Request.Context(), imageId); err != nil || !exist {
			NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("image", imageId))
			return
		}
	}

	// Create item
	var c []models.Color
	for _, colorId := range body.ColorsId {
		c = append(c, models.Color{Id: colorId})
	}
	var t []models.Tag
	for _, tag := range body.Tags {
		t = append(t, models.Tag{Name: tag})
	}
	var imgs []models.Image
	for _, imgId := range body.ImagesId {
		imgs = append(imgs, models.Image{Id: imgId})
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Get category and set
	category, err := h.services.Categories.GetById(ctx.Request.Context(), item.Category.Id)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	item.Category = category

	ctx.JSON(http.StatusOK, item)
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for i := range items {
		category, err := h.services.Categories.GetById(ctx.Request.Context(), items[i].Category.Id)
		if err != nil {
			NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		items[i].Category = category
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for i := range items {
		category, err := h.services.Categories.GetById(ctx.Request.Context(), items[i].Category.Id)
		if err != nil {
			NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		items[i].Category = category
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
// @Failure 400,404 {object} ErrorResponse
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	category, err := h.services.Categories.GetById(ctx.Request.Context(), item.Category.Id)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	item.Category = category

	ctx.JSON(http.StatusOK, item)
}

// @Summary Get items with category
// @Tags items-actions
// @Description get all items with provided category id
// @Accept json
// @Produce json
// @Param id path int true "category id"
// @Success 200 {array} models.Item
// @Failure 400 {object} ErrorResponse
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	category, err := h.services.Categories.GetById(ctx.Request.Context(), categoryId)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for i := range items {
		items[i].Category = category
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
// @Failure 400 {object} ErrorResponse
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	for i := range items {
		category, err := h.services.Categories.GetById(ctx.Request.Context(), items[i].Category.Id)
		if err != nil {
			NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		items[i].Category = category
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
// @Failure 400,404 {object} ErrorResponse
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
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	category, err := h.services.Categories.GetById(ctx.Request.Context(), item.Category.Id)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	item.Category = category

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
// @Failure 400 {object} ErrorResponse
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

// @Summary Update item
// @Security UsersAuth
// @Security AdminAuth
// @Tags items-actions
// @Description update item
// @Accept json
// @Produce json
// @Param id path string true "item id"
// @Param input body updateItemInput true "item body"
// @Success 200 {object} models.Item
// @Failure 400 {object} ErrorResponse
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

	// Check if item is existed
	if exist, err := h.services.Items.Exist(ctx.Request.Context(), itemId); err != nil || !exist {
		NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("item", itemId))
		return
	}

	//// Check if exist category
	//if exist, err := h.services.Categories.Exist(ctx.Request.Context(), body.CategoryId); err != nil || !exist {
	//	NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("category", body.CategoryId))
	//	return
	//}
	//
	//// Check if colors is existed
	//for _, colorId := range body.ColorsId {
	//	if exist, err := h.services.Colors.Exist(ctx.Request.Context(), colorId); err != nil || !exist {
	//		NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("color", colorId))
	//		return
	//	}
	//}

	// Check if at least one image
	if len(body.ImagesId) < 1 {
		err := errors.New("the item must contain at least 1 image")
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Check if images is existed
	for _, imageId := range body.ImagesId {
		if exist, err := h.services.Images.Exist(ctx.Request.Context(), imageId); err != nil || !exist {
			NewError(ctx, http.StatusNotFound, apperrors.ErrIdNotFound("image", imageId))
			return
		}
	}

	// Update item
	if err := h.services.Items.Update(ctx.Request.Context(), itemId, body.Name, body.Description,
		body.CategoryId, body.Tags, body.ColorsId, body.Price, body.Sku, body.ImagesId); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Return created item
	item, _ := h.services.Items.GetById(ctx.Request.Context(), itemId)

	// Get category and set
	category, err := h.services.Categories.GetById(ctx.Request.Context(), item.Category.Id)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	item.Category = category

	ctx.JSON(http.StatusOK, item)
}
