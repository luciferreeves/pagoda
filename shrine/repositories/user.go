package repositories

import (
	"shrine/database"
	"shrine/enums"
	"shrine/models"
	"time"
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

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func FindUserByVerification(hash string, verificationType enums.VerificationType) (*models.User, error) {
	var user models.User
	err := database.DB.Where("verification_hash = ? AND verification_type = ? AND verification_expiry > ?", hash, verificationType, time.Now()).First(&user).Error
	return &user, err
}