package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) corsMiddleware(c *gin.Context) {
	for _, domain := range h.cfg.HTTP.AllowedOrigins {
		c.Header("Access-Control-Allow-Origin", domain)
	}

	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
