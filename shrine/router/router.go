package router

import (
	"shrine/utils/shortcuts"
	"shrine/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	urls.Attach(router)
}

func ErrorHandler(context *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if fiberErr, ok := err.(*fiber.Error); ok {
		code = fiberErr.Code
	}

	switch code {
	case fiber.StatusBadRequest:
		return shortcuts.BadRequest(context, err)
	case fiber.StatusUnauthorized:
		return shortcuts.Unauthorized(context, err)
	case fiber.StatusForbidden:
		return shortcuts.Forbidden(context, err)
	case fiber.StatusNotFound:
		return shortcuts.NotFound(context, err)
	default:
		return shortcuts.InternalServerError(context, err)
	}
}
