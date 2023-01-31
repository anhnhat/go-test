package main

import (
	"fmt"
	"http-server/controllers"
	"http-server/initializers"
	"http-server/models"
	"http-server/routes"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	startTime           int64
	server              *gin.Engine
	userRouteController routes.UserRouteController
)

func init() {
	startTime = time.Now().Unix()

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&config)

	userRouteController = routes.CreateRouteUserController(
		controllers.CreateUserController(initializers.DB),
	)

	server = gin.Default()
}

func main() {
	router := server.Group("/api/v1")
	router.GET("/healthcheck", func(c *gin.Context) {
		uptime := int(time.Now().Unix() - startTime)

		c.Header("Cache-Control", "no-cache")
		c.JSON(200, gin.H{
			"status": "running",
			"uptime": strconv.Itoa(uptime) + " seconds",
		})
	})
	userRouteController.UserRoute(router)
	router.POST("/login", controllers.Login)
	router.POST("/signup", controllers.Signup)
	router.GET("/seed", func(ctx *gin.Context) {
		status := models.ProjectStatus{Name: "Planning"}
		role := models.Role{Name: "Member"}
		project := models.Project{Name: "Forth management", Description: "Software management system", StatusID: 1, Priority: "High"}

		res1 := initializers.DB.Create(&status)
		res2 := initializers.DB.Create(&role)
		res3 := initializers.DB.Create(&project)

		fmt.Println(project)
		fmt.Print(res1.Error)
		fmt.Print(res2.Error)
		fmt.Print(res3.Error)
	})
	router.GET("/projects", controllers.RetreiveProjects)
	router.POST("/project", controllers.CreateProject)
	server.Run()
}
