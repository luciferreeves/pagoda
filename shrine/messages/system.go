package messages

const (
	StorageNotConfigured    = "Storage not configured."
	InvalidServerPort       = "Invalid server port: %d."
	InvalidDatabaseDriver   = "Invalid database driver: %s."
	DSNCannotBeEmpty        = "Data source name (DSN) cannot be empty."
	ConfigMustBePointer     = "Config must be a pointer to struct."
	FailedLoadTemplate      = "Failed to load %s template: %v."
	CannotActionSelf        = "You cannot %s yourself."
	CannotActionOwner       = "You cannot %s the owner."
	OnlyOwnerCanActionAdmin = "Only the owner can %s an administrator."
	FailedLibSQLConnection  = "Failed to open libsql connection: %v."
)