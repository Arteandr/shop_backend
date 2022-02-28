package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		auth := users.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
		}
	}
}

type signUpInput struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// @Summary Sign-up
// @Tags users-auth
// @Description Users sign-up
// @Accept json
// @Produce json
// @Param input body signUpInput true "sign-up input"
// @Success 200 {string} string "ok"
// @Router /users/auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var body signUpInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

}
