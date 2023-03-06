package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"http-server/controllers"
	"http-server/repository"
)

func NewProjectRoutes(db *gorm.DB, r *gin.RouterGroup) {
	projectRepo := repository.NewProjectRepo(db)
	projectController := &controllers.ProjectController{
		ProjectUC: projectRepo,
	}

	routes := r.Group("/projects")
	routes.GET("", projectController.RetrieveProject)
	routes.GET("/:id", projectController.GetProjectByIdOrName)
	routes.POST("", projectController.CreateProject)
	routes.PUT("/:id", projectController.UpdateProject)
	routes.DELETE("/:id", projectController.DeleteProject)
}
