package initializers

import (
	"go_auth/models"
	"log"
)

func SyncDatabase() {
	// Sync the database with the models
	log.Println("Syncing database...")
	DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Permissions{},
		&models.GroupPermissions{},
		&models.UserPermissions{})
	log.Println("Database synced successfully")

	PopulateDatabase()
}

func PopulateDatabase() {
	// Populate the database with initial data
	log.Println("Populating database...")

	// Create initial users, groups, and permissions
	user := models.User{FirstName: "John", LastName: "Doe", UserName: "johndoe", Password: "password123"}
	group := models.Group{Name: "User"}
	permission := models.Permissions{Name: "Read"}
	permission2 := models.Permissions{Name: "Write"}

	DB.Create(&user)
	DB.Create(&group)
	DB.Create(&permission)
	DB.Create(&permission2)

	DB.Create(&models.GroupPermissions{GroupID: group.ID, PermissionID: permission.ID})
	DB.Create(&models.UserPermissions{UserID: user.ID, PermissionID: permission2.ID})
	DB.Create(&models.UserGroup{UserID: user.ID, GroupID: group.ID})

	log.Println("Database populated successfully")
}
