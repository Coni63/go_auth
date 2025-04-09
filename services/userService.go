package services

import (
	"errors"
	"fmt"
	"go_auth/initializers"
	"go_auth/models"

	"gorm.io/gorm"
)

func CreateUserFromRequest() (*models.User, error) {
	return nil, nil
}

func GetUserById(id string) (*models.User, error) {
	var user models.User
	result := initializers.DB.First(&user, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error // Return other errors normally
	}

	return &user, nil // Return found user
}

func GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	err := initializers.DB.Find(&users).Error
	return users, err
}

func UpdateUser(user *models.User) {
}

func DeleteUser(id string) error {
	result := initializers.DB.Delete(&models.User{}, id)

	// Check if the delete was successful
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	return nil
}
