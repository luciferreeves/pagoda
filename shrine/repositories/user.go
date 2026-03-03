package repositories

import (
	"shrine/database"
	"shrine/enums"
	"shrine/models"
)

func CreateUser(user *models.User) error {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	tx := database.DB
	if count == 0 {
		user.Role = enums.Admin
		tx = tx.Set("bypass_username_validation", true)
	}

	return tx.Create(user).Error
}

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}