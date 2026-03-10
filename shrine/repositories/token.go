package repositories

import (
	"shrine/database"
	"shrine/models"
	"time"
)

func CreateToken(token *models.Token) error {
	return database.DB.Create(token).Error
}

func FindValidToken(tokenHash string) (*models.Token, error) {
	var token models.Token
	err := database.DB.Preload("User").
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		First(&token).Error
	return &token, err
}

func DeleteToken(tokenHash string) error {
	return database.DB.Where("token_hash = ?", tokenHash).Delete(&models.Token{}).Error
}

func DeleteUserTokens(userID uint) {
	database.DB.Where("user_id = ?", userID).Delete(&models.Token{})
}