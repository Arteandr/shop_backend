package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/create [post]
func (h *Handler) createItem(ctx *gin.Context) {
	// Body binding
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if exist category
	if exist, err := h.services.Categories.Exist(body.CategoryId); err != nil || !exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong category id"})
		return
	}

	// Check if colors is existed
	for _, colorId := range body.ColorsId {
		if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong color[%d] id", colorId)})
			return
		}
	}

	// Check if at least one image
	if len(body.ImagesId) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "at least 1 image"})
		return
	}

	// Check if images is existed
	for _, imageId := range body.ImagesId {
		if exist, err := h.services.Images.Exist(imageId); err != nil || !exist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong image[%d] id", imageId)})
			return
		}
	}

	// Create item
	itemId, err := h.services.Items.Create(body.Name, body.Description, body.CategoryId, body.Sku, body.Price)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Link colors
	for i := 0; i < len(body.ColorsId); i++ {
		colorId := body.ColorsId[i]
		if err := h.services.Items.LinkColor(itemId, colorId); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	// Link tags if more than zero
	if len(body.Tags) > 0 {
		if err := h.services.Items.LinkTags(itemId, body.Tags); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	// Link images
	if err := h.services.Items.LinkImages(itemId, body.ImagesId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Return created item
	item, _ := h.services.Items.GetById(itemId)

	// Get category and set
	category, err := h.services.Categories.GetById(item.Category.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	items, err := h.services.Items.GetAll(sortOptions)

	for i := range items {
		category, err := h.services.Categories.GetById(items[i].Category.Id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
	items, err := h.services.Items.GetNew()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range items {
		category, err := h.services.Categories.GetById(items[i].Category.Id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	item, err := h.services.Items.GetById(itemId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("item %d not found", itemId)})
		return
	}

	category, err := h.services.Categories.GetById(item.Category.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	items, err := h.services.Items.GetByCategory(categoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	category, err := h.services.Categories.GetById(categoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong tag %s", tag)})
		return
	}
	items, err := h.services.Items.GetByTag(tag)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := range items {
		category, err := h.services.Categories.GetById(items[i].Category.Id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong sku"})
		return
	}

	item, err := h.services.Items.GetBySku(sku)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("item with sku %s not found", sku)})
		return
	}

	category, err := h.services.Categories.GetById(item.Category.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValid(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Items.Delete(ctx.Request.Context(), body.ItemsId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var body updateItemInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if item is existed
	if exist, err := h.services.Items.Exist(itemId); err != nil || !exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong item id %d", itemId)})
		return
	}

	// Check if exist category
	if exist, err := h.services.Categories.Exist(body.CategoryId); err != nil || !exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong category id"})
		return
	}

	// Check if colors is existed
	for _, colorId := range body.ColorsId {
		if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong color[%d] id", colorId)})
			return
		}
	}

	// Check if at least one image
	if len(body.ImagesId) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "at least 1 image"})
		return
	}

	// Check if images is existed
	for _, imageId := range body.ImagesId {
		if exist, err := h.services.Images.Exist(imageId); err != nil || !exist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong image[%d] id", imageId)})
			return
		}
	}

	// Update item
	if err := h.services.Items.Update(itemId, body.Name, body.Description,
		body.CategoryId, body.Tags, body.ColorsId, body.Price, body.Sku, body.ImagesId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Return created item
	item, _ := h.services.Items.GetById(itemId)

	// Get category and set
	category, err := h.services.Categories.GetById(item.Category.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	item.Category = category

	ctx.JSON(http.StatusOK, item)
}
