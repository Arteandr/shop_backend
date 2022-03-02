package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		auth := users.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
		}
	}
}

type authInput struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// @Summary Sign-up
// @Tags users-auth
// @Description Users sign-up
// @Accept json
// @Produce json
// @Param input body authInput true "sign-up input"
// @Success 200 {string} string "ok"
// @Router /users/auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var body authInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	exist := h.services.Users.EmailExist(strings.TrimSpace(body.Email))
	if exist {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "email already exist",
		})
		return
	}

	id, err := h.services.Users.SignUp(strings.TrimSpace(body.Email), strings.TrimSpace(body.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(ctx *gin.Context) {
	var body authInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokens, err := h.services.Users.SignIn(body.Email, body.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("refreshToken", tokens.RefreshToken, 2592000, "/", h.cfg.HTTP.Host, false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": tokens.AccessToken,
	})
}
