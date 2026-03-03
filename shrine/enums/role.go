package enums

type UserRole string

const (
	Member    UserRole = "member"
	Moderator UserRole = "moderator"
	Admin     UserRole = "admin"
)