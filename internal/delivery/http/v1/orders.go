package v1

import "github.com/gin-gonic/gin"

func (h *Handler) InitOrdersRoutes(api *gin.RouterGroup) {
  orders := api.Group("/orders")
  {
    orders.GET("/")
  }
}



