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

		authenticated := users.Group("/", h.userIdentify)
		{
			authenticated.GET("/me", h.getMe)
		}
	}
}

type authInput struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// @Summary Register
// @Tags users-auth
// @Description sign-up user
// @Accept json
// @Produce json
// @Param input body authInput true "email and password"
// @Success 200 {object} SignUpResponse
// @Failure 400,409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var body authInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	exist := h.services.Users.EmailExist(strings.TrimSpace(body.Email))
	if exist {
		ctx.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: "email already exist"})
		return
	}

	id, err := h.services.Users.SignUp(strings.TrimSpace(body.Email), strings.TrimSpace(body.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// @Summary Login
// @Tags users-auth
// @Description sign-in user to account
// @Accept json
// @Produce json
// @Param input body authInput true "email and password"
// @Success 200 {object} models.Tokens
// @Failure 400,404 {object} ErrorResponse
// @Router /users/auth/sign-in [post]
func (h *Handler) signIn(ctx *gin.Context) {
	var body authInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.services.Users.SignIn(body.Email, body.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.SetCookie("refreshToken", tokens.RefreshToken, 2592000, "/", h.cfg.HTTP.Host, false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": tokens.AccessToken,
	})
}

// @Summary Get current user
// @Security UsersAuth
// @Tags users-actions
// @Description get current user by auth token
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *Handler) getMe(ctx *gin.Context) {
	id, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.services.Users.GetUserById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
