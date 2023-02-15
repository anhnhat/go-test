package controllers

import (
	"fmt"
	"http-server/initializers"
	"http-server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Status = map[string]int{
	"Planning":   0,
	"Inprogress": 1,
	"Completed":  2,
	"Overdue":    3,
	"Cancel":     4,
}

func RetrieveProject(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var projects = []models.Project{}
	status, haveStatus := ctx.GetQuery("status")
	var query = initializers.DB.Limit(intLimit).Offset(offset)

	if haveStatus {
		query.Where("status_id = ?", Status[status])
	}
	query.Preload("Members").Find(&projects)

	if query.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": query.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projects,
	})
}

func CreateProject(ctx *gin.Context) {
	var payload *models.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	members := []*models.Member{}
	membersResult := initializers.DB.Where(payload.Members).Find(&members)

	if membersResult.Error != nil {
		fmt.Println(membersResult.Error)
	}

	project := models.Project{
		Name:           payload.Name,
		Description:    payload.Description,
		StatusID:       payload.StatusID,
		Priority:       payload.Priority,
		Health:         payload.Health,
		ClientName:     payload.ClientName,
		Members:        members,
		Budget:         payload.Budget,
		ActualReceived: payload.ActualReceived,
		StartDate:      payload.StartDate,
		EndDate:        payload.EndDate,
	}

	result := initializers.DB.Create(&project)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": result.Error,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Create project successfull",
		"data":    project,
	})
}

func GetProjectByIdOrName(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var project = models.Project{}
	id, err := strconv.Atoi(idParam)

	var result = initializers.DB
	if err != nil {
		result.Where("name LIKE ?", "%"+idParam+"%")
	} else {
		result.Where("id = ?", id)
	}
	result.Preload("Members").First(&project)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot find project",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload *models.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	members := []*models.Member{}
	membersResult := initializers.DB.Where(payload.Members).Find(&members)

	if membersResult.Error != nil {
		fmt.Println(membersResult.Error)
	}

	project := models.Project{
		Name:           payload.Name,
		Description:    payload.Description,
		StatusID:       payload.StatusID,
		Priority:       payload.Priority,
		Health:         payload.Health,
		Members:        members,
		StartDate:      payload.StartDate,
		EndDate:        payload.EndDate,
		ClientName:     payload.ClientName,
		Budget:         payload.Budget,
		ActualReceived: payload.ActualReceived,
	}

	result := initializers.DB.Model(&models.Project{}).Where("id = ?", id).Updates(&project)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot update project",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   project,
	})
}

func DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	result := initializers.DB.Delete(&models.Project{}, id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot delete project",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Deleted",
	})
}
