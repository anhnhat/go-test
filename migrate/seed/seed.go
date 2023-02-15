package main

import (
	"http-server/controllers"
	"http-server/initializers"
	"http-server/models"
	"log"
	"os"
	"time"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	initializers.ConnectDB(&config)
}

func main() {
	options := os.Args
	if len(options) < 2 {
		log.Fatal("🚀 Missing option: init || reset")
	}

	seedType := options[1]
	switch seedType {
	case "init":
		{
			seedData()
			log.Println("Seed data successfull")
		}
	case "reset":
		{
			initializers.DB.Exec("DROP SCHEMA public CASCADE;")
			initializers.DB.Exec("CREATE SCHEMA public;")
			log.Println("All table deleted, run migration to recreate !")
		}
	case "update":
		{
			// TODO
			log.Println("Update table !")
		}
	}
}

func seedData() {
	hashedPassword, _ := controllers.HashPassword("123456")
	user := models.User{
		Name:       "Admin",
		Email:      "admin@gmail.com",
		Password:   hashedPassword,
		AllowLogin: true,
	}
	role := []models.Role{
		{Name: "Developer"},
		{Name: "Desinger"},
		{Name: "QA/QC"},
		{Name: "PM"},
	}
	member := models.Member{
		Name:  "Anh Nhat",
		Email: "anhnhat@gmail.com",
		Team: models.Team{
			Name: "Team 1",
		},
	}
	project := models.Project{
		Name:        "Forth management",
		Description: "Super super man",
		StatusID:    1,
		Priority:    "High",
		Health: models.Health{
			Health:       "Strong",
			HealthReason: "Nothing",
		},
		ClientName:     "TFH",
		Budget:         2000.0,
		ActualReceived: 2200.0,
		StartDate:      time.Now(),
		EndDate:        time.Now(),
	}

	initializers.DB.Create(&user)
	initializers.DB.Create(&role)
	initializers.DB.Create(&member)
	initializers.DB.Create(&project)
}
