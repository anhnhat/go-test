package controllers

import (
	"http-server/initializers"
	"http-server/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func Signup(ctx *gin.Context) {
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

	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

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

func Login(ctx *gin.Context) {
	var payload *LoginUserRequest

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", payload.Email)
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), err
}
