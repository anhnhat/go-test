package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"http-server/controllers"
	"http-server/repository"
)

func NewMemberRoute(db *gorm.DB, r *gin.RouterGroup) {
	memberRepo := repository.NewMemberRepo(db)
	memberController := &controllers.MemberController{
		MemberUC: memberRepo,
	}

	routes := r.Group("/members")
	routes.GET("", memberController.RetreiveMembers)
	routes.GET("/:id", memberController.GetMember)
	routes.POST("", memberController.CreateMember)
	routes.PUT("/:id", memberController.UpdateMember)
	routes.DELETE("/:id", memberController.DeleteMember)
}
