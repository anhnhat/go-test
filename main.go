package main

import (
	"http-server/initializers"
	"http-server/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	startTime int64
	ginServer *gin.Engine
)

func init() {
	startTime = time.Now().Unix()
	log.SetFormatter(&log.JSONFormatter{})

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&config)

	ginServer = gin.Default()
	ginServer.Use(cors.Default())
}

func main() {
	routes.Setup(initializers.DB, ginServer)
	ginServer.Run()
}
