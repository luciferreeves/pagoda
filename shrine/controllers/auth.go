package controllers

import (
	"errors"
	"shrine/enums"
	"shrine/models"
	"shrine/repositories"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/emails"
	"shrine/utils/logger"
	"shrine/utils/meta"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterController(context *fiber.Ctx) error {
	body, err := meta.Body[types.RegisterRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("Invalid request body."))
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

	token, err := auth.GenerateToken()
	if err != nil {
		return InternalServerError(context, errors.New("Failed to generate verification token."))
	}

	user.SetVerification(auth.HashToken(token), time.Now().Add(24*time.Hour), enums.Activation)

	if err := repositories.UpdateUser(&user); err != nil {
		return InternalServerError(context, errors.New("Failed to store verification token."))
	}

	if err := emails.SendActivation(user.Email, user.Username, token); err != nil {
		logger.Errorf("Auth", "Failed to send verification email to %s: %v", user.Email, err)
	}

	return Created(context, types.MessageResponse{
		Message: "Your account has been created. Please check your email to verify your account.",
	})
}

func LoginController(context *fiber.Ctx) error {
	body, err := meta.Body[types.LoginRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("Invalid request body."))
	}

	user, err := repositories.FindUserByUsername(body.Username)
	if err != nil {
		return Unauthorized(context, errors.New("Invalid username or password."))
	}

	if !user.CanAuthenticate() {
		return Forbidden(context, errors.New("Your account has been banned or disabled."))
	}

	if !user.CheckPassword(body.Password) {
		return Unauthorized(context, errors.New("Invalid username or password."))
	}

	if !user.IsVerified() {
		return Forbidden(context, errors.New("Your email address has not been verified. Please check your inbox."))
	}

	token, err := auth.IssueToken(context, user.ID)
	if err != nil {
		return InternalServerError(context, errors.New("Failed to create session."))
	}

	return Success(context, types.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func VerifyController(context *fiber.Ctx) error {
	body, err := meta.Body[types.VerifyRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("Invalid request body."))
	}

	if body.Token == "" {
		return BadRequest(context, errors.New("Verification token is required."))
	}

	tokenHash := auth.HashToken(body.Token)

	user, err := repositories.FindUserByVerification(tokenHash, enums.Activation)
	if err != nil {
		return BadRequest(context, errors.New("Your verification link is invalid or has expired."))
	}

	user.VerifyEmail()

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to verify your account."))
	}

	return Success(context, types.MessageResponse{
		Message: "Your email has been verified successfully. You can now log in.",
	})
}

func LogoutController(context *fiber.Ctx) error {
	tokenHash := auth.GetTokenHash(context)
	if err := repositories.DeleteToken(tokenHash); err != nil {
		return InternalServerError(context, errors.New("Failed to end your session."))
	}

	return Success(context, types.MessageResponse{
		Message: "You have been logged out successfully.",
	})
}

func MeController(context *fiber.Ctx) error {
	user := auth.GetUser(context)
	return Success(context, user.ToResponse())
}