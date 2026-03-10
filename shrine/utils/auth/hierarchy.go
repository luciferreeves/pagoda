package auth

import (
	"fmt"
	"shrine/messages"
	"shrine/models"
)

func ValidateHierarchy(admin *models.User, target *models.User, action string) error {
	if target.ID == admin.ID {
		return fmt.Errorf(messages.CannotActionSelf, action)
	}

	if target.IsOwner() {
		return fmt.Errorf(messages.CannotActionOwner, action)
	}

	if target.IsAdmin() && !admin.IsOwner() {
		return fmt.Errorf(messages.OnlyOwnerCanActionAdmin, action)
	}

	return nil
}