package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
	}
}

type userSignUpInput struct {
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/sign-up [post]
func (h *Handler) userSignUp(ctx *gin.Context) {
	var body userSignUpInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.services.Users.SignUp(ctx.Request.Context(), body.Email, body.Login, body.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
