package routes

import (
	"http-server/controllers"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func CreateRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (userRouteController *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("", userRouteController.userController.RetreiveUsers)
	router.POST("", userRouteController.userController.CreateNewUser)
}
