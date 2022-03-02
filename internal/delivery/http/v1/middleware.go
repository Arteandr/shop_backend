package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_backend/internal/models"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"

	userCtx = "userId"
)

func (h *Handler) userIdentify(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Set(userCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(authorizationHeader)
	if authHeader == "" {
		return "", models.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", models.ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", models.ErrEmptyAuthHeader
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getIdByContext(c *gin.Context, context string) (int, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, errors.New(context + " not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New(context + " is of invalid type")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New(context + " convert failed")
	}

	return id, nil
}
