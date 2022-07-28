package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func (h *Handler) InitFilesRoutes(api *gin.RouterGroup) {
	assets := api.Group("/assets")
	{
		assets.Static("/", "./assets")
		assets.POST("/upload", h.uploadFile)
	}
}

func (h *Handler) uploadFile(ctx *gin.Context) {
	photo, err := ctx.FormFile("photo")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	filename := filepath.Base(photo.Filename)
	if err := ctx.SaveUploadedFile(photo, "./assets/"+filename); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Images.Upload(filename); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
