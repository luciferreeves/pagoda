package auth

import (
	"fmt"
	"shrine/models"
)

func ValidateHierarchy(admin *models.User, target *models.User, action string) error {
	if target.ID == admin.ID {
		return fmt.Errorf("You cannot %s yourself.", action)
	}

	if target.IsOwner() {
		return fmt.Errorf("You cannot %s the owner.", action)
	}

	if target.IsAdmin() && !admin.IsOwner() {
		return fmt.Errorf("Only the owner can %s an administrator.", action)
	}

	return nil
}