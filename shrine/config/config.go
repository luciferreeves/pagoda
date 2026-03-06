package config

import (
	"shrine/utils/env"
	"shrine/utils/logger"

	"github.com/joho/godotenv"
)

var (
	Server   server
	Database database
	SMTP     smtp
	Storage  storage
)

func init() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		logger.Infof("Config", "No .env file found. Environment variables will be used directly.")
	}

	if err := env.Parse(&Server); err != nil {
		logger.Fatalf("Config", "Failed to parse server config: %v", err)
	}

	if err := env.Parse(&Database); err != nil {
		logger.Fatalf("Config", "Failed to parse database config: %v", err)
	}

	if err := env.Parse(&SMTP); err != nil {
		logger.Fatalf("Config", "Failed to parse SMTP config: %v", err)
	}

	if err := env.Parse(&Storage); err != nil {
		logger.Fatalf("Config", "Failed to parse storage config: %v", err)
	}

	if Server.Debug {
		logger.SetDebug(true)
	}

	if err := verifyConfig(); err != nil {
		logger.Fatalf("Config", "Configuration verification failed: %v", err)
	}

	logger.Successf("Config", "Configuration loaded successfully")
}
