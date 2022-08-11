package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitImagesRoutes(api *gin.RouterGroup) {
	images := api.Group("/images")
	{
		admins := images.Group("/", h.userIdentity, h.adminIdentify)
		{
			admins.POST("/", h.uploadFile)
			admins.GET("/", h.getAllImages)
			admins.DELETE("/", h.deleteImages)
		}

	}
}

// @Summary Upload image
// @Security UsersAuth
// @Security AdminAuth
// @Tags images-actions
// @Description upload image
// @Accept json
// @Produce json
// @Param photo formData file true "photo to upload"
// @Success 200 {object} IdResponse
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

	ctx.JSON(http.StatusOK, IdResponse{Id: id})
}

// @Summary Get all images
// @Security UsersAuth
// @Security AdminAuth
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

type deleteImagesInput struct {
	ImagesId []int `json:"imagesId" binding:"required"`
}

func (i *deleteImagesInput) isValid() error {
	if len(i.ImagesId) < 1 {
		return fmt.Errorf("wrong images id length %d", len(i.ImagesId))
	}

	return nil
}

// @Summary Delete images
// @Security UsersAuth
// @Security AdminAuth
// @Tags images-actions
// @Description delete images by id
// @Accept json
// @Produce json
// @Param input body deleteImagesInput true "images id info"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /images/ [delete]
func (h *Handler) deleteImages(ctx *gin.Context) {
	var body deleteImagesInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValid(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Images.Delete(ctx.Request.Context(), body.ImagesId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
