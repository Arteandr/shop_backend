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
		items.POST("/create", h.createItem)
		items.GET("/:id", h.getItemById)
		items.GET("/sku/:sku", h.getItemBySku)
		items.GET("/category/:id", h.getItemsByCategory)
		items.GET("/tag/:name", h.getItemsByTag)
		items.DELETE("/:id", h.deleteItem)

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

	ctx.JSON(http.StatusOK, item)
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
	tag := ctx.Param("id")
	if len(tag) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong tag"})
		return
	}

	items, err := h.services.Items.GetByTag(tag)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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

	ctx.JSON(http.StatusOK, item)
}

// @Summary Delete item
// @Tags items-actions
// @Description delete item by id
// @Accept json
// @Produce json
// @Param id path int true "item id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /items/{id} [delete]
func (h *Handler) deleteItem(ctx *gin.Context) {
	strItemId := ctx.Param("id")
	itemId, err := strconv.Atoi(strItemId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err = h.services.Items.Delete(itemId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
