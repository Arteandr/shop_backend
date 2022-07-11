package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitItemsRoutes(api *gin.RouterGroup) {
	items := api.Group("/items")
	{
		items.POST("/create", h.createItem)
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

	itemId, err := h.services.Items.Create(body.Name, body.Description, body.CategoryId, body.Tags, body.Sku)
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

	ctx.Status(http.StatusOK)
}
