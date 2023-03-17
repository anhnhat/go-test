package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAuthService interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthService struct {
	DB *gorm.DB
}

// func NewAuthService(DB *gorm.DB) IAuthService {
// 	return AuthService{
// 		DB: DB,
// 	}
// }

// func ()
