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

	}
}

type createItemInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CategoryId  int      `json:"categoryId" binding:"required"`
	Tags        []string `json:"tags,omitempty"`
	ColorsId    []int    `json:"colors" binding:"required"`
	Sku         string   `json:"sku" binding:"required"`
}

func (h *Handler) createItem(ctx *gin.Context) {
	var body createItemInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if exist, err := h.services.Categories.Exist(body.CategoryId); err != nil || !exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong category id"})
		return
	}

	for _, colorId := range body.ColorsId {
		if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("wrong color[%d] id", colorId)})
			return
		}
	}

	itemId, err := h.services.Items.Create(body.Name, body.Description, body.CategoryId, body.Sku)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i := 0; i < len(body.ColorsId); i++ {
		colorId := body.ColorsId[i]
		if err := h.services.Items.LinkColor(itemId, colorId); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	if len(body.Tags) > 0 {
		if err := h.services.Items.LinkTags(itemId, body.Tags); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	ctx.Status(http.StatusOK)
}

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
