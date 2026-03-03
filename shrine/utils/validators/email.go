package validators

import "strings"

func IsValidEmail(email string) bool {
	if len(email) < 5 || len(email) > 255 {
		return false
	}

	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return false
	}

	dotIndex := strings.LastIndex(email, ".")
	if dotIndex < atIndex+2 || dotIndex >= len(email)-1 {
		return false
	}

	return true
}