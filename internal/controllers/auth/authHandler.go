package auth

import (
	"net/http"
	"strconv"
	"time"

	"http-server/cmd/server/config"
	"http-server/internal/appctx"
	"http-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type authHandler struct {
	appCtx *appctx.AppCtx
}

func NewAuthHandler(appCtx *appctx.AppCtx) authHandler {
	return authHandler{
		appCtx: appCtx,
	}
}

func (ah *authHandler) Signup(ctx *gin.Context) {
	var payload *models.CreateUserRequest

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	now := time.Now()
	newUser := models.User{
		ID:         uuid.New(),
		Name:       payload.Name,
		Email:      payload.Email,
		Password:   string(hashedPassword),
		AllowLogin: true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// result := ac.DB.Create(&newUser)
	// if result.Error != nil {
	// 	ctx.JSON(http.StatusBadGateway, gin.H{
	// 		"status":  "error",
	// 		"message": result.Error.Error(),
	// 	})
	// 	return
	// }

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   userResponse,
	})
}

func (ah *authHandler) Login(ctx *gin.Context) {
	var payload *LoginUserRequest
	config, _ := config.LoadConfig(".")
	jwtService := ah.appCtx.GetJwtService()

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	result := ah.appCtx.DB.First(&user, "email = ?", payload.Email)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})
		return
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if compareErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	if !user.AllowLogin {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "You don't have permission to login",
		})
		return
	}

	tokenExpHour, _ := strconv.Atoi(config.JWT_EXPIRY_HOUR)
	refreshTokenExpHour, _ := strconv.Atoi(config.JWT_REFRESH_TOKEN_EXPIRY_HOUR)
	tokenExp := time.Now().Add(time.Hour * time.Duration(tokenExpHour)).Unix()
	refreshTokenExp := time.Now().Add(time.Hour * time.Duration(refreshTokenExpHour)).Unix()

	token, err := jwtService.GenerateToken(user.ID.String(), tokenExp)
	refreshToken, err := jwtService.GenerateToken(user.ID.String(), refreshTokenExp)
	ctx.SetCookie("token", token, int(tokenExpHour*3600), "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, int(refreshTokenExp*3600), "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), err
}

func FetchMe(ctx *gin.Context) {}
