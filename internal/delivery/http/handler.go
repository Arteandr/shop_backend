package delivery

import (
	"github.com/gin-gonic/gin"
	"shop_backend/internal/config"
	v1 "shop_backend/internal/delivery/http/v1"
	"shop_backend/internal/service"
)

type Handler struct {
	services *service.Services
	cfg      *config.Config
}

func NewHandler(services *service.Services, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	h.InitApi(r)

	return r
}

func (h *Handler) InitApi(r *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.cfg)
	api := r.Group("/api")
	{
		handlerV1.Init(api)
	}
}