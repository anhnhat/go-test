package project

import (
	"fmt"
	"http-server/internal/appctx"
	"http-server/internal/models"
	"http-server/internal/services/project"
	"http-server/internal/services/timesheet"
	"time"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type projectHandler struct {
	appCtx           *appctx.AppCtx
	projectService   project.IProjectService
	timesheetService timesheet.ITimeSheet
}

func NewProjectHandler(appCtx *appctx.AppCtx) projectHandler {
	return projectHandler{
		appCtx:           appCtx,
		projectService:   appCtx.GetProjectService(),
		timesheetService: appCtx.GetTimesheetService(),
	}
}

func (ph *projectHandler) RetrieveProject(ctx *gin.Context) {
	projects, err := ph.projectService.GetAll(ctx)

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

func (ph *projectHandler) CreateProject(ctx *gin.Context) {
	var payload *models.CreateProjectRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	project, err := ph.projectService.Create(payload)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Create timesheet when assign member to project
	if len(payload.Members) != 0 {
		for _, memberId := range payload.Members {
			timesheet := models.TimeSheet{
				MemberID:  int(memberId),
				ProjectID: int(project.ID),
				DayLog: models.DayLog{
					Date: datatypes.Date(time.Now()),
				},
			}
			if err := ph.timesheetService.Create(timesheet); err != nil {
				fmt.Println("Cannot create timesheet for project")
			}
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Create project successfull",
		"data":    project,
	})
}

func (ph *projectHandler) GetProjectByIdOrName(ctx *gin.Context) {
	id := ctx.Param("id")

	project, err := ph.projectService.GetByIdOrName(id)
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

func (ph *projectHandler) UpdateProject(ctx *gin.Context) {
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

	project, err := ph.projectService.Update(idInt, payload, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "Cannot update project",
		})
		return
	}

	if err := ph.projectService.AssignMembersToProject(idInt, payload.Members); err != nil {
		println("Cannot save to association member_projects")
	}

	// Create timesheet for each project-member
	if len(payload.Members) != 0 {
		for _, memberId := range payload.Members {
			timesheet := models.TimeSheet{
				MemberID:  int(memberId),
				ProjectID: idInt,
				DayLog: models.DayLog{
					Date: datatypes.Date(time.Now()),
				},
			}
			if err := ph.timesheetService.Create(timesheet); err != nil {
				fmt.Println("Cannot create timesheet for project")
			}
		}
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

func (ph *projectHandler) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)

	// TODO: Delete all relation

	err := ph.projectService.Delete(idInt)
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
