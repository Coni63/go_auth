package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	// Connect to the database using the DSN from environment variables
	// The DSN should be in the format: "host=localhost user=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	// You can set this in your .env file or directly in your environment variables
	log.Println("Connecting to database...")
	var err error
	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("DSN not set in environment variables")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Connected to database successfully")
}
