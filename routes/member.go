package routes

import (
	"http-server/controllers"

	"github.com/gin-gonic/gin"
)

func MemberRoutes(r *gin.RouterGroup) {
	routes := r.Group("/members")
	routes.GET("", controllers.RetreiveMembers)
	routes.GET("/:id", controllers.GetMember)
	routes.POST("", controllers.CreateMember)
	// routes.PUT("/:id", controllers.UpdateProject)
	// routes.DELETE("/:id", controllers.DeleteProject)
}
