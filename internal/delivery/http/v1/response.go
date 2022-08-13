package v1

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type IdResponse struct {
	Id int `json:"id"`
}

func NewError(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(code, ErrorResponse{Error: err.Error()})
}
