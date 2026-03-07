package services

import (
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/hypertext"
	"shrine/utils/crypto"
	"shrine/utils/emails"
	"shrine/utils/logger"
	"time"
)

func SendVerification(citizen *models.User, verificationType enums.VerificationType) *hypertext.ServiceError {
	token, err := crypto.GenerateToken()
	if err != nil {
		return fail(enums.Internal, messages.FailedGenerateToken)
	}

	citizen.SetVerification(crypto.HashToken(token), time.Now().Add(24*time.Hour), verificationType)

	if err := repositories.UpdateUser(citizen); err != nil {
		return fail(enums.Internal, messages.FailedStoreToken)
	}

	switch verificationType {
	case enums.Activation:
		if err := emails.SendActivation(citizen.Email, citizen.Username, token); err != nil {
			logger.Errorf("Auth", "Failed to send activation email to %s: %v", citizen.Email, err)
		}
	}

	return nil
}