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

func (h *Handler) userIdentity(ctx *gin.Context) {
	id, err := h.parseAuthHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userCtx, id)
}

func (h *Handler) parseAuthHeader(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		return "", models.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", models.ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getIdByContext(ctx *gin.Context, context string) (int, error) {
	idFromCtx, ok := ctx.Get(context)
	if !ok {
		return 0, errors.New(context + " not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New(context + " is of invalid type")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New(context + " is of invalid type")
	}

	return id, nil
}
