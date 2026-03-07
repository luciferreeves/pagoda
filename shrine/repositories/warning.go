package repositories

import (
	"shrine/database"
	"shrine/models"
	"shrine/utils/crypto"
	"shrine/utils/meta"

	"gorm.io/gorm"
)

func CreateWarning(adminID uint, userID uint, title string, message string) (*models.Warning, error) {
	var warning models.Warning
	ref := crypto.SystemRef()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		letter := models.Letter{
			Title:     title,
			IsSystem:  true,
			SystemRef: ref,
		}
		if err := tx.Create(&letter).Error; err != nil {
			return err
		}

		participant := models.LetterParticipant{
			LetterID: letter.ID,
			UserID:   userID,
		}
		if err := tx.Create(&participant).Error; err != nil {
			return err
		}

		msg := models.LetterMessage{
			LetterID: letter.ID,
			Body:     message,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		warning = models.Warning{
			UserID:    userID,
			AdminID:   adminID,
			LetterID:  letter.ID,
			SystemRef: ref,
			Message:   message,
		}
		if err := tx.Create(&warning).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("warning_count", gorm.Expr("warning_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	return &warning, err
}

func DeactivateWarning(warning *models.Warning) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		warning.Active = false
		if err := tx.Save(warning).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ? AND warning_count > 0", warning.UserID).Update("warning_count", gorm.Expr("warning_count - 1")).Error; err != nil {
			return err
		}

		return nil
	})
}

func ListWarningsForUser(userID uint, p meta.Pagination) ([]models.Warning, int64) {
	var warnings []models.Warning
	var total int64

	query := database.DB.Model(&models.Warning{}).Where("user_id = ?", userID)
	query.Count(&total)
	p.Apply(query.Order("created_at desc")).Preload("Admin").Find(&warnings)

	return warnings, total
}

func FindWarningByRef(ref string) (*models.Warning, error) {
	var warning models.Warning
	err := database.DB.Preload("Admin").Where("system_ref = ?", ref).First(&warning).Error
	return &warning, err
}