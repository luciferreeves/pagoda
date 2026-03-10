package enums

type DatabaseDriver string

const (
	SQLite   DatabaseDriver = "sqlite"
	Postgres DatabaseDriver = "postgres"
	LibSQL   DatabaseDriver = "libsql"
)
