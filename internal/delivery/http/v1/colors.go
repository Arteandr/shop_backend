package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	apperrors "shop_backend/pkg/errors"
	"strconv"
)

func (h *Handler) InitColorsRoutes(api *gin.RouterGroup) {
	colors := api.Group("/colors")
	{
		admins := colors.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/all/:id", h.addColorToItems)
			admins.DELETE("/all/:id", h.deleteColorFromItems)
			admins.POST("/create", h.createColor)
			admins.PUT("/:id", h.updateColor)
			admins.DELETE("/", h.deleteColors)
		}

		colors.GET("/", h.getAllColors)
		colors.GET("/:id", h.getColorById)
	}
}

type createColorInput struct {
	Name  string   `json:"name" binding:"required"`
	Hex   string   `json:"hex" binding:"required"`
	Price *float64 `json:"price" binding:"required"`
}

// @Summary Create a new color
// @Security UsersAuth
// @Security AdminAuth
// @Tags colors-actions
// @Description create a new color
// @Accept json
// @Produce json
// @Param input body models.Color true "input body"
// @Success 200 {object} IdResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/create [post]
func (h *Handler) createColor(ctx *gin.Context) {
	var color createColorInput
	if err := ctx.BindJSON(&color); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	colorId, err := h.services.Colors.Create(ctx.Request.Context(), color.Name, color.Hex, *color.Price)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, IdResponse{Id: colorId})
}

type updateColorInput struct {
	Name  string   `json:"name" binding:"required"`
	Hex   string   `json:"hex" binding:"required"`
	Price *float64 `json:"price" binding:"required"`
}

// @Summary Update color
// @Security UsersAuth
// @Security AdminAuth
// @Tags colors-actions
// @Description update color
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Param input body updateColorInput true "input body"
// @Success 200 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/{id} [put]
func (h *Handler) updateColor(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	var color updateColorInput
	if err := ctx.BindJSON(&color); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	exist, err := h.services.Colors.Exist(ctx.Request.Context(), colorId)
	if !exist {
		err = errors.New("wrong color id")
		NewError(ctx, http.StatusNotFound, err)
		return
	} else if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := h.services.Colors.Update(ctx.Request.Context(), colorId, color.Name, color.Hex, *color.Price); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type deleteColorsInput struct {
	ColorsId []int `json:"colorsId" binding:"required"`
}

func (i *deleteColorsInput) isValid() error {
	if len(i.ColorsId) < 1 {
		return fmt.Errorf("wrong images id length %d", len(i.ColorsId))
	}

	return nil
}

// @Summary Delete colors
// @Security UsersAuth
// @Security AdminAuth
// @Tags colors-actions
// @Description delete color by id
// @Accept json
// @Produce json
// @Param input body deleteColorsInput true "images id info"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors [delete]
func (h *Handler) deleteColors(ctx *gin.Context) {
	var body deleteColorsInput
	if err := ctx.BindJSON(&body); err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidBody)
		return
	}

	if err := body.isValid(); err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Colors.Delete(ctx.Request.Context(), body.ColorsId); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Delete color from all items
// @Security UsersAuth
// @Security AdminAuth
// @Tags colors-actions
// @Description delete color by id from all items
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/all/{id} [delete]
func (h *Handler) deleteColorFromItems(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	if err := h.services.Colors.DeleteFromItems(ctx.Request.Context(), colorId); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Add color to all items
// @Security UsersAuth
// @Security AdminAuth
// @Tags colors-actions
// @Description Add color by id to all items
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Success 200 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/all/{id} [post]
func (h *Handler) addColorToItems(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	exist, err := h.services.Colors.Exist(ctx.Request.Context(), colorId)
	if !exist {
		err = errors.New("wrong color id")
		NewError(ctx, http.StatusNotFound, err)
		return
	} else if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := h.services.Colors.AddToItems(ctx.Request.Context(), colorId); err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get color by id
// @Tags colors-actions
// @Description get color by id
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Success 200 ""
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/{id} [get]
func (h *Handler) getColorById(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		NewError(ctx, http.StatusBadRequest, apperrors.ErrInvalidParam)
		return
	}

	exist, err := h.services.Colors.Exist(ctx.Request.Context(), colorId)
	if !exist {
		err = errors.New("wrong color id")
		NewError(ctx, http.StatusNotFound, err)
		return
	} else if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	color, err := h.services.Colors.GetById(ctx.Request.Context(), colorId)
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, color)
}

// @Summary Get all colors
// @Tags colors-actions
// @Description get all colors
// @Accept json
// @Produce json
// @Success 200 {array} models.Color
// @Failure 500 {object} ErrorResponse
// @Router /colors [get]
func (h *Handler) getAllColors(ctx *gin.Context) {
	colors, err := h.services.Colors.GetAll(ctx.Request.Context())
	if err != nil {
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, colors)
}
