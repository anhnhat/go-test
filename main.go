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
	server.Use(CORSMiddleware())
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
