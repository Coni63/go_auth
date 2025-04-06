package initializers

import (
	"go_auth/models"
	"log"
)

func SyncDatabase() {
	// Sync the database with the models
	log.Println("Syncing database...")
	DB.AutoMigrate(&models.User{})
	log.Println("Database synced successfully")

	// for i := 0; i < 100; i++ {
	// 	DB.Create(&models.User{Name: "John", Lastname: "Doe"})
	// }
}
