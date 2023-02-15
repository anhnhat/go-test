package routes

import (
	"http-server/controllers"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(r *gin.RouterGroup) {
	routes := r.Group("/projects")
	routes.GET("", controllers.RetrieveProject)
	routes.GET("/:id", controllers.GetProjectByIdOrName)
	routes.POST("/", controllers.CreateProject)
	routes.PUT("/:id", controllers.UpdateProject)
	routes.DELETE("/:id", controllers.DeleteProject)
}
