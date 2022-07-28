package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitImagesRoutes(api *gin.RouterGroup) {
	assets := api.Group("/assets")
	{
		assets.POST("/upload", h.uploadFile)
		assets.GET("/images", h.getAllImages)
	}
}

func (h *Handler) uploadFile(ctx *gin.Context) {
	photo, err := ctx.FormFile("photo")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.services.Images.Upload(photo)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) getAllImages(ctx *gin.Context) {
	images, err := h.services.Images.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	for i, image := range images {
		images[i].Filename = "/files/" + image.Filename
	}

	ctx.JSON(http.StatusOK, images)
}
