package routes

import (
	"http-server/controllers"
	"http-server/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTimeSheetRoutes(db *gorm.DB, r *gin.RouterGroup) {
	timeSheetRepo := repository.NewTimeSheetRepo(db)
	timeSheetController := &controllers.TimeSheetController{
		TimeSheetUC: timeSheetRepo,
	}

	routes := r.Group("/timesheet")
	routes.GET("", timeSheetController.GetAll)
	routes.GET("/:id", timeSheetController.GetByMemberId)
	routes.POST("", timeSheetController.Create)
}
