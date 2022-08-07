package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"regexp"
	"shop_backend/internal/models"
	"strings"
)

func (h *Handler) InitUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/refresh", h.userRefresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/logout", h.userLogout)
			authenticated.GET("/me", h.userGetMe)

			authenticated.PUT("/email", h.userUpdateEmail)
			authenticated.PUT("/password", h.userUpdatePassword)
			authenticated.PUT("/info", h.userUpdateInfo)
			authenticated.PUT("/address", h.userUpdateAddress)
		}
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

// @Summary User sign-up
// @Tags users-auth
// @Description create user account
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {object} models.User
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

	ctx.SetCookie("refresh_token", tokens.RefreshToken, 2592000, "/", "localhost", false, true)

	// Hide refresh token
	tokens.RefreshToken = ""
	ctx.JSON(http.StatusOK, tokens)
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Tokens
// @Failure 400,404,401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/refresh [post]
func (h *Handler) userRefresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.services.Users.RefreshTokens(ctx.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.SetCookie("refresh_token", tokens.RefreshToken, 2592000, "/", "localhost", false, true)

	// Hide refresh token
	tokens.RefreshToken = ""
	ctx.JSON(http.StatusOK, tokens)
}

// @Summary Get current user
// @Security UsersAuth
// @Tags users-auth
// @Description get current user by authentication header
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *Handler) userGetMe(ctx *gin.Context) {
	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.services.Users.GetMe(ctx.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Logout current user
// @Security UsersAuth
// @Tags users-auth
// @Description logout current user by authentication header
// @Accept  json
// @Produce  json
// @Success 200 ""
// @Failure 500 {object} ErrorResponse
// @Router /users/logout [post]
func (h *Handler) userLogout(ctx *gin.Context) {
	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Users.Logout(ctx.Request.Context(), userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

type userUpdateEmailInput struct {
	Email string `json:"email" binding:"required"`
}

func (u *userUpdateEmailInput) isValidEmail() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("wrong email")
	}

	const emailLength = 30
	if len(u.Email) > emailLength {
		return errors.New(fmt.Sprintf("email length must not exceed %d characters", emailLength))
	}

	return nil
}

// @Summary User update email
// @Tags users-auth
// @Description update current user email
// @Accept  json
// @Produce  json
// @Param input body userUpdateEmailInput true "email info"
// @Success 200 ""
// @Failure 400, 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/email [put]
func (h *Handler) userUpdateEmail(ctx *gin.Context) {
	var body userUpdateEmailInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValidEmail(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Users.UpdateEmail(ctx.Request.Context(), userId, body.Email); err != nil {
		switch err.(type) {
		case models.ErrUniqueValue:
			ctx.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	ctx.Status(http.StatusOK)
}

type userUpdatePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (u *userUpdatePasswordInput) isValidPassword() error {
	if len(u.NewPassword) < 6 || len(u.NewPassword) > 16 {
		return errors.New("wrong new password length")
	}

	return nil
}

// @Summary User update password
// @Tags users-auth
// @Description update current user password
// @Accept  json
// @Produce  json
// @Param input body userUpdatePasswordInput true "password info"
// @Success 200 ""
// @Failure 400, 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/password [put]
func (h *Handler) userUpdatePassword(ctx *gin.Context) {
	var body userUpdatePasswordInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValidPassword(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Users.UpdatePassword(ctx, userId, body.OldPassword, body.NewPassword); err != nil {
		if errors.Is(err, models.ErrOldPassword) {
			ctx.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

type userUpdateInfoInput struct {
	Login       string `json:"login" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneCode   string `json:"phoneCode" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

func (u *userUpdateInfoInput) isValidInfo() error {
	// Check login
	if len(u.Login) < 2 || len(u.Login) > 15 {
		return errors.New("wrong login length")
	}
	const loginPattern = `^[A-Za-z0-9]+$`
	if matched, _ := regexp.MatchString(loginPattern, u.Login); !matched {
		return errors.New("wrong login")
	}
	// Check first name
	if len(u.FirstName) < 1 || len(u.FirstName) > 20 {
		return errors.New("wrong first name length")
	}
	// Check last name
	if len(u.FirstName) < 1 || len(u.FirstName) > 20 {
		return errors.New("wrong first name length")
	}
	// Check phone code
	if len(u.PhoneCode) < 1 || len(u.PhoneCode) > 5 {
		return errors.New("wrong phone code length")
	}
	// Check phone number
	if len(u.PhoneNumber) < 1 || len(u.PhoneNumber) > 15 {
		return errors.New("wrong phone number length")
	}

	return nil
}

// @Summary User update info
// @Tags users-auth
// @Description update current user info
// @Accept  json
// @Produce  json
// @Param input body userUpdateInfoInput true "info body"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/info [put]
func (h *Handler) userUpdateInfo(ctx *gin.Context) {
	var body userUpdateInfoInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := body.isValidInfo(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Users.UpdateInfo(ctx, userId, body.Login, body.FirstName, body.LastName, body.PhoneCode, body.PhoneNumber); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

type userUpdateAddressInput struct {
	InvoiceAddress  models.Address `json:"invoiceAddress" binding:"required"`
	ShippingAddress models.Address `json:"shippingAddress" binding:"required"`
}

func (u *userUpdateAddressInput) isDiffer() bool {
	if strings.TrimSpace(u.InvoiceAddress.Country) == strings.TrimSpace(u.ShippingAddress.Country) &&
		strings.TrimSpace(u.InvoiceAddress.City) == strings.TrimSpace(u.ShippingAddress.City) &&
		strings.TrimSpace(u.InvoiceAddress.Street) == strings.TrimSpace(u.ShippingAddress.Street) &&
		u.InvoiceAddress.Zip == u.ShippingAddress.Zip {
		return false
	} else {
		return true
	}
}

// @Summary User update address
// @Tags users-auth
// @Description update current user address
// @Accept  json
// @Produce  json
// @Param input body userUpdateAddressInput true "address info"
// @Success 200 ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/info [put]
func (h *Handler) userUpdateAddress(ctx *gin.Context) {
	var body userUpdateAddressInput
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	different := body.isDiffer()

	userId, err := getIdByContext(ctx, userCtx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.services.Users.UpdateAddress(ctx, userId, different, body.InvoiceAddress, body.ShippingAddress); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
