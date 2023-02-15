package routes

import (
	"http-server/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
}
