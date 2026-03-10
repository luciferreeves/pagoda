package database

import (
	"shrine/config"
	"shrine/enums"
	"shrine/utils/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	MaxOpenConnections    = 25
	MaxIdleConnections    = 5
	ConnectionMaxLifetime = time.Hour
)

var DB *gorm.DB

func init() {
	var dialector gorm.Dialector

	switch enums.DatabaseDriver(config.Database.Driver) {
	case enums.SQLite:
		dialector = sqlite.Open(config.Database.DSN)
	case enums.Postgres:
		dialector = postgres.Open(config.Database.DSN)
	default:
		logger.Fatalf("Database", "Invalid database driver: %s", config.Database.Driver)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})

	if err != nil {
		logger.Fatalf("Database", "Error connecting to database: %v", err)
	}

	db, err := DB.DB()
	if err != nil {
		logger.Fatalf("Database", "Failed to get underlying sql.DB: %v", err)
	}
	db.SetMaxOpenConns(MaxOpenConnections)
	db.SetMaxIdleConns(MaxIdleConnections)
	db.SetConnMaxLifetime(ConnectionMaxLifetime)

	logger.Successf("Database", "Database connection established successfully")

	migrate()
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
