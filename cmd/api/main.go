package main

import (
	"log"

	"github.com/TIA-PARTNERS-GROUP/tia-api/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	config := configs.LoadConfig()

	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	print(db)

	log.Println("Database connection successful.")
}
