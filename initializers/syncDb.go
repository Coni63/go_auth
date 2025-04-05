package initializers

import (
	"go_auth/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})

	// for i := 0; i < 100; i++ {
	// 	DB.Create(&models.User{Name: "John", Lastname: "Doe"})
	// }
}
