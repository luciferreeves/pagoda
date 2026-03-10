package controllers

import (
	"shrine/services"
	"shrine/types/account"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func RegisterController(context *fiber.Ctx) error {
	body, err := meta.Body[account.RegisterRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.Register(body, meta.Request(context).IP)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func LoginController(context *fiber.Ctx) error {
	body, err := meta.Body[account.LoginRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	citizen, serviceErr := services.Authenticate(body, meta.Request(context).IP)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	token, err := auth.IssueToken(context, citizen.ID)
	if err != nil {
		return shortcuts.InternalServerError(context, err)
	}

	return shortcuts.Success(context, account.AuthResponse{
		Token: token,
		User:  citizen.ToResponse(),
	})
}

func VerifyController(context *fiber.Ctx) error {
	body, err := meta.Body[account.VerifyRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.VerifyAccount(body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func ResendActivationController(context *fiber.Ctx) error {
	body, err := meta.Body[account.ResendActivationRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.ResendActivation(body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func LogoutController(context *fiber.Ctx) error {
	tokenHash := auth.GetTokenHash(context)

	result, serviceErr := services.RevokeToken(tokenHash)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func MeController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	return shortcuts.Success(context, citizen.ToResponse())
}

func HeartbeatController(context *fiber.Ctx) error {
	return shortcuts.NoContent(context)
}