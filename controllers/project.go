package controllers

import (
	"fmt"
	"http-server/initializers"
	"http-server/models"

	"github.com/gin-gonic/gin"
)

func RetreiveProjects(ctx *gin.Context) {
	var projects = []models.Project{}

	result := initializers.DB.Find(&projects)

	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"status":  "error",
			"message": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data":   projects,
	})

	fmt.Println(projects)
}

func CreateProject(ctx *gin.Context) {
	var project models.Project

	err := ctx.ShouldBindJSON(&project)

	if err != nil {

	}
}
