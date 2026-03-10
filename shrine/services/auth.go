package services

import (
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/account"
	"shrine/types/common"
	"shrine/types/hypertext"
	"shrine/utils/crypto"
)

func Register(request account.RegisterRequest, ip string) (*common.MessageResponse, *hypertext.ServiceError) {
	if repositories.IsIPBanned(ip) {
		return nil, fail(enums.Forbidden, messages.IPBanned)
	}

	citizen := models.User{
		Username:    request.Username,
		Email:       request.Email,
		DisplayName: request.DisplayName,
		IP:          ip,
	}

	if err := citizen.SetPassword(request.Password); err != nil {
		return nil, fail(enums.BadRequest, err.Error())
	}

	if err := repositories.CreateUser(&citizen); err != nil {
		return nil, mapRegistrationError(err)
	}

	if serviceErr := SendVerification(&citizen, enums.Activation); serviceErr != nil {
		return nil, serviceErr
	}

	return &common.MessageResponse{Message: messages.AccountCreated}, nil
}

func Authenticate(request account.LoginRequest, ip string) (*models.User, *hypertext.ServiceError) {
	citizen, err := repositories.FindUserByUsername(request.Username)
	if err != nil {
		return nil, fail(enums.Unauthorized, messages.InvalidUsernameOrPassword)
	}

	if !citizen.IsStaff() && repositories.IsIPBanned(ip) {
		return nil, fail(enums.Forbidden, messages.IPBanned)
	}

	if citizen.ClearExpiredDisable() {
		repositories.UpdateUser(citizen)
	}

	if !citizen.CanAuthenticate() {
		return nil, fail(enums.Forbidden, messages.AccountBannedOrDisabled)
	}

	if !citizen.CheckPassword(request.Password) {
		return nil, fail(enums.Unauthorized, messages.InvalidUsernameOrPassword)
	}

	if !citizen.IsVerified() {
		return nil, fail(enums.Forbidden, messages.EmailNotVerified)
	}

	citizen.IP = ip
	repositories.UpdateUser(citizen)

	return citizen, nil
}

func VerifyAccount(request account.VerifyRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	if request.Token == "" {
		return nil, fail(enums.BadRequest, messages.VerificationTokenRequired)
	}

	verificationType := enums.VerificationType(request.Type)
	citizen, err := repositories.FindUserByVerification(crypto.HashToken(request.Token), verificationType)
	if err != nil {
		return nil, fail(enums.BadRequest, messages.VerificationLinkInvalid)
	}

	switch verificationType {
	case enums.Activation:
		citizen.VerifyEmail()
	default:
		return nil, fail(enums.BadRequest, messages.InvalidVerificationType)
	}

	if err := repositories.UpdateUser(citizen); err != nil {
		return nil, fail(enums.Internal, messages.FailedVerifyAccount)
	}

	return &common.MessageResponse{Message: messages.EmailVerified}, nil
}

func ResendActivation(request account.ResendActivationRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	citizen, err := repositories.FindUserByEmail(request.Email)
	if err != nil {
		return nil, fail(enums.BadRequest, messages.NoAccountWithEmail)
	}

	if citizen.IsVerified() {
		return nil, fail(enums.BadRequest, messages.AccountAlreadyVerified)
	}

	if serviceErr := SendVerification(citizen, enums.Activation); serviceErr != nil {
		return nil, serviceErr
	}

	return &common.MessageResponse{Message: messages.VerificationEmailSent}, nil
}

func RevokeToken(tokenHash string) (*common.MessageResponse, *hypertext.ServiceError) {
	if err := repositories.DeleteToken(tokenHash); err != nil {
		return nil, fail(enums.Internal, messages.FailedEndSession)
	}
	return &common.MessageResponse{Message: messages.LoggedOut}, nil
}