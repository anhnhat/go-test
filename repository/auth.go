package repository

import "github.com/gin-gonic/gin"

type IAuthRepo interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}
