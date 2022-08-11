package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "shop_backend/docs"
	"shop_backend/internal/config"
	v1 "shop_backend/internal/delivery/http/v1"
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

func (h *Handler) Init() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.InitApi(r)

	return r
}

func (h *Handler) InitApi(r *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.cfg, h.tokenManager)
	api := r.Group("/")
	{
		handlerV1.Init(api)
	}
}
