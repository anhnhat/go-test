package main

import (
	"log"
	"os"
	"time"

	"http-server/cmd/server/config"
	"http-server/infras/db"
	controllers "http-server/internal/controllers/auth"
	"http-server/internal/models"
	"http-server/internal/services/member"

	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	memberService member.MemberService
)

func init() {
	configInstance, err := config.LoadConfig(".")
	if err != nil {
		panic("Cannot load env")
	}

	dsn := config.GetDsn(&configInstance)
	DB, err = db.NewDB(dsn)
	if err != nil {
		panic("Cannot connect to database")
	}
}

func main() {
	options := os.Args
	if len(options) < 2 {
		log.Fatal("🚀 Missing option: drop || migrate || seed")
	}

	seedType := options[1]
	switch seedType {
	case "seed":
		{
			seedData()
			log.Println("👍 Seed data successfull")
		}
	case "migrate":
		{
			DB.AutoMigrate(
				&models.User{},
				&models.Project{},
				&models.Role{},
				&models.Team{},
				&models.Member{},
				&models.TimeSheet{},
				&models.TimeSheetSegment{},
			)
			log.Println("👍 Migration complete!!!")
		}
	case "drop":
		{
			DB.Exec("DROP SCHEMA public CASCADE;")
			DB.Exec("CREATE SCHEMA public;")
			log.Println("👍 All table deleted, run migrate to recreate !")
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
	roles := []models.Role{
		{Name: "Developer"},
		{Name: "Desinger"},
		{Name: "QA/QC"},
		{Name: "PM"},
	}
	teams := []models.Team{
		{Name: "QA/QC"},
		{Name: "Technical"},
		{Name: "Designer"},
	}
	members := []models.Member{
		{
			Name:      "Anh Nhat",
			Email:     "anhnhat@gmail.com",
			TeamID:    1,
			WorkModel: 0,
			Salary:    2.1,
			OtherCost: 2.2,
			StartDate: time.Now(),
		},
		{
			Name:      "Anh Khoa",
			Email:     "anh@gmail.com",
			TeamID:    2,
			WorkModel: 1,
			Salary:    2.1,
			OtherCost: 2.2,
			StartDate: time.Now(),
		},
		{
			Name:      "Be",
			Email:     "be@gmail.com",
			TeamID:    1,
			WorkModel: 2,
			Salary:    2.1,
			OtherCost: 2.2,
			StartDate: time.Now(),
		},
		{
			Name:      "Ha Ngo",
			Email:     "hango@gmail.com",
			TeamID:    2,
			WorkModel: 1,
			Salary:    2000,
			OtherCost: 210,
			StartDate: time.Now(),
		},
		{
			Name:      "Huynh Dang",
			Email:     "huynh@gmail.com",
			TeamID:    1,
			WorkModel: 2,
			Salary:    3100,
			OtherCost: 110,
			StartDate: time.Now(),
		},
	}
	projects := []models.Project{
		{
			Name:        "Forth management",
			Description: "Project management system",
			StatusID:    1,
			Priority:    "High",
			Health: models.Health{
				Health:       "Strong",
				HealthReason: "By some unexpected reason",
			},
			ClientName:     "TFH",
			Budget:         2000.0,
			ActualReceived: 2200.0,
			StartDate:      time.Now(),
			EndDate:        time.Now(),
		},
		{
			Name:        "Enjoy NFT Land",
			Description: "NFT Marketplace",
			StatusID:    2,
			Priority:    "High",
			Health: models.Health{
				Health:       "Weak",
				HealthReason: "By some unexpected reason",
			},
			ClientName:     "TFH",
			Budget:         2000.0,
			ActualReceived: 2200.0,
			StartDate:      time.Now(),
			EndDate:        time.Now(),
		},
		{
			Name:        "Circle",
			Description: "Co-working admin system",
			StatusID:    3,
			Priority:    "High",
			Health: models.Health{
				Health:       "Weak",
				HealthReason: "By some unexpected reason",
			},
			ClientName:     "TFH",
			Budget:         2000.0,
			ActualReceived: 2200.0,
			StartDate:      time.Now(),
			EndDate:        time.Now(),
		},
		{
			Name:        "B52 wine",
			Description: "Wine checker",
			StatusID:    3,
			Priority:    "High",
			Health: models.Health{
				Health:       "Weak",
				HealthReason: "By some unexpected reason",
			},
			ClientName:     "Me",
			Budget:         2000.0,
			ActualReceived: 2200.0,
			StartDate:      time.Now(),
			EndDate:        time.Now(),
		},
	}

	DB.Create(&user)
	DB.Create(&roles)
	DB.Create(&teams)
	DB.Create(&members)
	DB.Create(&projects)

	memberService = member.MemberService{DB: DB}
	memberService.AssignProjectsToMember(1, []int{1, 2})
	memberService.AssignProjectsToMember(2, []int{2, 3})
	memberService.AssignRolesToMember(1, []int{1})
	memberService.AssignRolesToMember(2, []int{2})
}
