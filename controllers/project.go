package controllers

import (
	"net/http"
	"strconv"

	"http-server/models"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	ProjectUC models.IProjectRepo
}

func (pc *ProjectController) RetrieveProject(ctx *gin.Context) {
	projects, err := pc.ProjectUC.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projects,
	})
}

func (pc *ProjectController) CreateProject(ctx *gin.Context) {
	var payload *models.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	project, err := pc.ProjectUC.Create(payload)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Create project successfull",
		"data":    project,
	})
}

func (pc *ProjectController) GetProjectByIdOrName(ctx *gin.Context) {
	id := ctx.Param("id")

	project, err := pc.ProjectUC.GetByIdOrName(id)
	if err != nil {
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

func (pc *ProjectController) UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	var payload *models.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	project, err := pc.ProjectUC.Update(idInt, payload, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot update project",
		})
		return
	}

	if err := pc.ProjectUC.AssignMembersToProject(idInt, payload.Members); err != nil {
		println("Cannot save to association member_projects")
	}

	if err != nil {
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

func (pc *ProjectController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	err := pc.ProjectUC.Delete(idInt)
	if err != nil {
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
