package routes

import (
	"net/http"

	"http-server/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, server *gin.Engine) {
	routes := server.Group("/api/v1")
	// TODO: refactor auth routes
	AuthRoutes(routes)

	private := routes.Group("/")
	private.Use(services.MiddleWare())
	private.GET("/me", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"access_token": "skjdf",
			"status":       "success",
		})
		return
	})
	NewMemberRoute(db, routes)
	NewProjectRoutes(db, routes)
	NewTimeSheetRoutes(db, routes)
}
