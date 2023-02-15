package main

import (
	"http-server/controllers"
	"http-server/initializers"
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
	// server.Use(cors.Default())
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
	routes.ProjectRoutes(router)
	routes.AuthRoutes(router)
	routes.MemberRoutes(router)

	server.Run()
}
