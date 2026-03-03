package middleware

import (
	"shrine/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func request() fiber.Handler {
	return func(context *fiber.Ctx) error {
		req := meta.BuildRequest(context)
		context.Locals(meta.RequestKey, req)
		return context.Next()
	}
}
