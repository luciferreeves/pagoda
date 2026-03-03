package auth

import "github.com/gofiber/fiber/v2"

func IsAuthenticated(context *fiber.Ctx) bool {
	// We will implement token based authentication in the future
	return true
}

func RequireAuthentication(handler fiber.Handler) fiber.Handler {
	return func(context *fiber.Ctx) error {
		if !IsAuthenticated(context) {
			return fiber.ErrUnauthorized
		}
		return handler(context)
	}
}
