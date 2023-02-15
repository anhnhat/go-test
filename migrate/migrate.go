package main

import (
	"fmt"
	"http-server/initializers"
	"http-server/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Role{},
		&models.Team{},
		&models.Member{},
	)
	fmt.Println("👍 Migration complete!!!")
}
