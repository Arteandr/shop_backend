package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"regexp"
	"shop_backend/internal/models"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
	}
}

type userSignUpInput struct {
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *userSignUpInput) isValidEmail() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("wrong email")
	}

	const emailLength = 30
	if len(u.Email) > emailLength {
		return errors.New(fmt.Sprintf("email length must not exceed %d characters", emailLength))
	}

	return nil
}

func (u *userSignUpInput) isValidLogin() error {
	if len(u.Login) < 2 || len(u.Login) > 15 {
		return errors.New("wrong login length")
	}

	// Include all latin alphabet and numbers 0-9
	const loginPattern = `^[A-Za-z0-9]+$`
	if matched, _ := regexp.MatchString(loginPattern, u.Login); !matched {
		return errors.New("wrong login")
	}

	return nil
}

func (u *userSignUpInput) isValidPassword() error {
	if len(u.Password) < 6 || len(u.Password) > 16 {
		return errors.New("wrong password length")
	}

	return nil
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
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

type userSignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *userSignInInput) loginIsEmail() bool {
	if _, err := mail.ParseAddress(u.Login); err == nil {
		return true
	} else {
		return false
	}
}

// @Summary User SignIn
// @Tags users-auth
// @Description login into user account
// @Accept  json
// @Produce  json
// @Param input body userSignInInput true "sign in info"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/sign-in [post]
func (h *Handler) userSignIn(ctx *gin.Context) {
	var body userSignInInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var findBy string
	if body.loginIsEmail() {
		findBy = "email"
	} else {
		findBy = "login"
	}

	tokens, err := h.services.Users.SignIn(ctx.Request.Context(), findBy, body.Login, body.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
