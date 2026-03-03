package database

import (
	"shrine/utils/logger"
)

func migrate() {
	err := DB.AutoMigrate(
	// Models will be added here as they are created
	)
	if err != nil {
		logger.Fatalf("Database", "Error during database migration: %v", err)
	}

	logger.Successf("Database", "Database migration completed successfully")
}
