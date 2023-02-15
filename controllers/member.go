package controllers

import (
	"http-server/initializers"
	"http-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetreiveMembers(ctx *gin.Context) {
	members := []models.Member{}

	result := initializers.DB.Find(&members)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error,
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": members,
	})
}

func CreateMember(ctx *gin.Context) {
	var payload *models.CreateMemberRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	roles := []*models.Role{}
	projects := []*models.Project{}

	initializers.DB.Find(&roles, payload.Roles)
	initializers.DB.Find(&projects, payload.Projects)

	member := models.Member{
		Name:      payload.Name,
		Email:     payload.Email,
		Roles:     roles,
		Projects:  projects,
		TeamID:    int(payload.TeamID),
		WorkModel: models.WorkModel(payload.WorkModel),
		Salary:    float32(payload.Salary),
		OtherCost: float32(payload.OtherCost),
	}

	result := initializers.DB.Create(&member)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot create member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}

func GetMember(ctx *gin.Context) {
	id := ctx.Param("id")

	member := models.Member{}
	result := initializers.DB.First(&member, id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot find member",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   member,
	})
}
