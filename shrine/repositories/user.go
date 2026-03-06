package repositories

import (
	"shrine/database"
	"shrine/enums"
	"shrine/models"
	"shrine/utils/meta"
	"time"
)

func CreateUser(user *models.User) error {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	tx := database.DB
	if count == 0 {
		user.Role = enums.Owner
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

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindUserByVerification(hash string, verificationType enums.VerificationType) (*models.User, error) {
	var user models.User
	err := database.DB.Where("verification_hash = ? AND verification_type = ? AND verification_expiry > ?", hash, verificationType, time.Now()).First(&user).Error
	return &user, err
}

func UpdateLastSeen(user *models.User) {
	database.DB.Model(user).Update("last_seen_at", user.LastSeenAt)
}

func CountCitizens() int64 {
	var count int64
	database.DB.Model(&models.User{}).Where("email_verified = ?", true).Count(&count)
	return count
}

func CountOnline() int64 {
	var count int64
	database.DB.Model(&models.User{}).Where("last_seen_at > ?", time.Now().Add(-5*time.Minute)).Count(&count)
	return count
}

func NewestCitizens(limit int) []models.User {
	var users []models.User
	database.DB.Where("email_verified = ?", true).Order("created_at desc").Limit(limit).Find(&users)
	return users
}

func OnlineCitizens(limit int) []models.User {
	var users []models.User
	database.DB.Where("last_seen_at > ?", time.Now().Add(-5*time.Minute)).Order("last_seen_at desc").Limit(limit).Find(&users)
	return users
}

func ListUsers(p meta.Pagination, search string) ([]models.User, int64) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("username LIKE ? OR display_name LIKE ? OR email LIKE ?", like, like, like)
	}

	query.Count(&total)
	p.Apply(query.Order("created_at desc")).Find(&users)

	return users, total
}