package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"regexp"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in")
	}
}

type userSignUpInput struct {
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *userSignUpInput) isValidEmail() error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("wrong email")
	}

	const emailLength = 30
	if len(user.Email) > emailLength {
		return errors.New(fmt.Sprintf("email length must not exceed %d characters", emailLength))
	}

	return nil
}

func (user *userSignUpInput) isValidLogin() error {
	if len(user.Login) < 2 || len(user.Login) > 15 {
		return errors.New("wrong login length")
	}

	// Include all latin alphabet and numbers 0-9
	const loginPattern = `^[A-Za-z0-9]+$`
	if matched, _ := regexp.MatchString(loginPattern, user.Login); !matched {
		return errors.New("wrong login")
	}

	return nil
}

func (user *userSignUpInput) isValidPassword() error {
	if len(user.Password) < 6 || len(user.Password) > 16 {
		return errors.New("wrong password length")
	}

	return nil
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

	if err := body.isValidEmail(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValidLogin(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValidPassword(); err != nil {
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
