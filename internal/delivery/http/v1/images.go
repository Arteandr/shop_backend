package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) InitImagesRoutes(api *gin.RouterGroup) {
	assets := api.Group("/images")
	{
		assets.POST("/", h.uploadFile)
		assets.GET("/", h.getAllImages)
		assets.DELETE("/:id", h.deleteImage)
	}
}

// @Summary Upload image
// @Tags images-actions
// @Description upload image
// @Accept json
// @Produce json
// @Param photo formData file true "photo to upload"
// @Success 200 {object} UploadFileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /images/ [post]
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

	ctx.JSON(http.StatusOK, UploadFileResponse{Id: id})
}

// @Summary Get all images
// @Tags images-actions
// @Description get all images
// @Accept json
// @Produce json
// @Success 200 {array} models.Image
// @Failure 500 {object} ErrorResponse
// @Router /images/ [get]
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

// @Summary Delete image
// @Tags images-actions
// @Description delete image by id
// @Accept json
// @Produce json
// @Param id path int true "image id"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /images/{id} [delete]
func (h *Handler) deleteImage(ctx *gin.Context) {
	strImageId := ctx.Param("id")
	imageId, err := strconv.Atoi(strImageId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Images.Delete(imageId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
