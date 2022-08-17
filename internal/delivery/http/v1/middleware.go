package v1

import (
	"errors"
	"net/http"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"

	userCtx    = "userId"
	optionsCtx = "options"

	ASC  = "ASC"
	DESC = "DESC"
)

func (h *Handler) sort(defaultSortField, defaultSortOrder string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sortBy, ok := ctx.GetQuery("sort_by")
		if !ok {
			sortBy = defaultSortField
		}

		sortOrder, ok := ctx.GetQuery("sort_order")
		if !ok {
			sortOrder = defaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != ASC && upperSortOrder != DESC {
				ctx.Status(http.StatusBadRequest)
				return
			}
		}

		options := models.SortOptions{
			Field: sortBy,
			Order: sortOrder,
		}

		ctx.Set(optionsCtx, options)
	}
}

func getSortOptions(ctx *gin.Context) (models.SortOptions, error) {
	optionsFromCtx, ok := ctx.Get(optionsCtx)
	if !ok {
		return models.SortOptions{}, errors.New(optionsCtx + " not found")
	}

	options, ok := optionsFromCtx.(models.SortOptions)
	if !ok {
		return models.SortOptions{}, errors.New(optionsCtx + " is of invalid type")
	}

	return options, nil
}

func (h *Handler) completedIdentify(ctx *gin.Context) {
	id, err := getIdByContext(ctx, userCtx)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	completed, err := h.services.Users.IsCompleted(ctx, id)
	if !completed {
		NewError(ctx, http.StatusForbidden, apperrors.ErrUserNotCompleted)
		return
	}
}

func (h *Handler) userIdentity(ctx *gin.Context) {
	id, err := h.parseAuthHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userCtx, id)
}

func (h *Handler) adminIdentify(ctx *gin.Context) {
	id, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.services.Users.GetMe(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	if user.Admin != true {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
}

func (h *Handler) parseAuthHeader(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		return "", apperrors.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", apperrors.ErrInvalidAuthHeader
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
