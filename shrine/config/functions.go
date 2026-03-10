package config

import (
	"fmt"
	"shrine/enums"
	"shrine/messages"
)

func verifyConfig() error {
	if Server.Port <= 0 || Server.Port > 65535 {
		return fmt.Errorf(messages.InvalidServerPort, Server.Port)
	}

	if !verifyDatabaseDriver(enums.DatabaseDriver(Database.Driver)) {
		return fmt.Errorf(messages.InvalidDatabaseDriver, Database.Driver)
	}

	if Database.DSN == "" {
		return fmt.Errorf(messages.DSNCannotBeEmpty)
	}

	return nil
}

func verifyDatabaseDriver(driver enums.DatabaseDriver) bool {
	switch driver {
	case enums.SQLite, enums.Postgres:
		return true
	default:
		return false
	}
}
