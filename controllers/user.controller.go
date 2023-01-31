package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"http-server/models"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	DB *gorm.DB
}

func CreateUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (userController *UserController) RetreiveUsers(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := userController.DB.Limit(intLimit).Offset(offset).Find(&users)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": results.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"results": len(users),
		"data":    users,
	})
}

func (userController *UserController) CreateNewUser(ctx *gin.Context) {
	var payload *models.CreateUserRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	now := time.Now()
	newUser := models.User{
		ID:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := userController.DB.Create(&newUser)
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
