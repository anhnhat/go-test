package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (DB *gorm.DB, err error) {
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("Failed to connect to the Database")
		return nil, err
	}
	fmt.Println("🚀 Connected Successfully to the Database")
	return DB, nil
}
