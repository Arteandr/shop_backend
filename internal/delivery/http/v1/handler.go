package v1

import (
	"github.com/gin-gonic/gin"
	"shop_backend/internal/config"
	"shop_backend/internal/service"
)

type Handler struct {
	services *service.Services
	cfg      *config.Config
}

func NewHandler(services *service.Services, cfg *config.Config) *Handler {
	return &Handler{services: services, cfg: cfg}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.InitUsersRoutes(v1)
	}
}
