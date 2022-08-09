package v1

import (
	"github.com/gin-gonic/gin"
	"shop_backend/internal/config"
	"shop_backend/internal/service"
	"shop_backend/pkg/auth"
)

type Handler struct {
	services     *service.Services
	cfg          *config.Config
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, cfg *config.Config, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		cfg:          cfg,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.InitUsersRoutes(v1)
		h.InitItemsRoutes(v1)
		h.InitColorsRoutes(v1)
		h.InitCategoriesRoutes(v1)
		h.InitImagesRoutes(v1)
		h.InitDeliveryRoutes(v1)
	}
}
