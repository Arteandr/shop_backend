package v1

import "github.com/gin-gonic/gin"

func (h *Handler) InitItemsRoutes(api *gin.RouterGroup) {
	items := api.Group("/items")
	{
		items.POST("/create", h.createItem)
	}
}

func (h *Handler) createItem(ctx *gin.Context) {

}
