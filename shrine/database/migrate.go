package database

import (
	"shrine/models"
	"shrine/utils/logger"
)

func migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Letter{},
		&models.LetterParticipant{},
		&models.LetterMessage{},
		&models.LetterAttachment{},
		&models.Warning{},
		&models.AuditLog{},
		&models.TicketCategory{},
		&models.Ticket{},
		&models.TicketMessage{},
	)
	if err != nil {
		logger.Fatalf("Database", "Error during database migration: %v", err)
	}

	logger.Successf("Database", "Database migration completed successfully")
}
