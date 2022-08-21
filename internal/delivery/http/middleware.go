package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) corsMiddleware(c *gin.Context) {
	protocol := "http"
	if h.cfg.HTTP.HTTPS {
		protocol = "https"
	}

	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", fmt.Sprintf("%s://%s", protocol, h.cfg.HTTP.Host))
	c.Header("Access-Control-Allow-Origin", fmt.Sprintf("%s://admin.%s", protocol, h.cfg.HTTP.Host))
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
