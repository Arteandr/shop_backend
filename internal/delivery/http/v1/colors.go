package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) InitColorsRoutes(api *gin.RouterGroup) {
	colors := api.Group("/colors")
	{
		colors.GET("/", h.getAllColors)
		colors.GET("/:id", h.getColorById)
		colors.POST("/create", h.createColor)
		colors.POST("/all/:id", h.addColorToItems)
		colors.DELETE("/all/:id", h.deleteColorFromItems)
		colors.DELETE("/:id", h.deleteColor)
		colors.PUT("/:id", h.updateColor)
	}
}

type createColor struct {
	Id    int      `json:"id,omitempty"`
	Name  string   `json:"name" binding:"required"`
	Hex   string   `json:"hex" binding:"required"`
	Price *float64 `json:"price" binding:"required"`
}

// @Summary Create a new color
// @Tags colors-actions
// @Description create a new color
// @Accept json
// @Produce json
// @Param input body models.Color true "input body"
// @Success 200 {object} CreateColorResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/create [post]
func (h *Handler) createColor(ctx *gin.Context) {
	var color createColor
	if err := ctx.BindJSON(&color); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	colorId, err := h.services.Colors.Create(color.Name, color.Hex, *color.Price)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, CreateColorResult{ColorId: colorId})
}

type updateColorInput struct {
	Name  string   `json:"name" binding:"required" db:"name"`
	Hex   string   `json:"hex" binding:"required" db:"hex"`
	Price *float64 `json:"price" binding:"required" db:"price"`
}

// @Summary Update color
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var color updateColorInput
	if err := ctx.BindJSON(&color); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "wrong input body"})
		return
	}

	if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
		if !exist {
			err = errors.New("wrong color id")
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Colors.Update(colorId, color.Name, color.Hex, *color.Price); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Delete colors
// @Tags colors-actions
// @Description delete color by id
// @Accept json
// @Produce json
// @Param id path int true "color id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /colors/{id} [delete]
func (h *Handler) deleteColor(ctx *gin.Context) {
	strColorId := ctx.Param("id")
	colorId, err := strconv.Atoi(strColorId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Colors.Delete(colorId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Delete color from all items
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Colors.DeleteFromItems(colorId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Add color to all items
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
		if !exist {
			err = errors.New("wrong color id")
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Colors.AddToItems(colorId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if exist, err := h.services.Colors.Exist(colorId); err != nil || !exist {
		if !exist {
			err = errors.New("wrong color id")
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	color, err := h.services.Colors.GetById(colorId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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
// @Router /colors/ [get]
func (h *Handler) getAllColors(ctx *gin.Context) {
	colors, err := h.services.Colors.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, colors)
}
