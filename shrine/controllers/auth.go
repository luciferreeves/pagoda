package controllers

import (
	"errors"
	"shrine/models"
	"shrine/repositories"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func RegisterController(context *fiber.Ctx) error {
	body, err := meta.Body[types.RegisterRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("invalid request body"))
	}

	user := models.User{
		Username:    body.Username,
		Email:       body.Email,
		DisplayName: body.DisplayName,
	}

	if err := user.SetPassword(body.Password); err != nil {
		return BadRequest(context, err)
	}

	if err := repositories.CreateUser(&user); err != nil {
		return BadRequest(context, err)
	}

	token, err := auth.IssueToken(context, user.ID)
	if err != nil {
		return InternalServerError(context, errors.New("failed to create session"))
	}

	return Created(context, types.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func LoginController(context *fiber.Ctx) error {
	body, err := meta.Body[types.LoginRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("invalid request body"))
	}

	user, err := repositories.FindUserByUsername(body.Username)
	if err != nil {
		return Unauthorized(context, errors.New("invalid credentials"))
	}

	if !user.CanAuthenticate() {
		return Forbidden(context, errors.New("account is not eligible for authentication"))
	}

	if !user.CheckPassword(body.Password) {
		return Unauthorized(context, errors.New("invalid credentials"))
	}

	token, err := auth.IssueToken(context, user.ID)
	if err != nil {
		return InternalServerError(context, errors.New("failed to create session"))
	}

	return Success(context, types.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func LogoutController(context *fiber.Ctx) error {
	tokenHash := auth.GetTokenHash(context)
	if err := repositories.DeleteToken(tokenHash); err != nil {
		return InternalServerError(context, errors.New("failed to end session"))
	}

	return Success(context, types.MessageResponse{
		Message: "logged out",
	})
}

func MeController(context *fiber.Ctx) error {
	user := auth.GetUser(context)
	return Success(context, user.ToResponse())
}