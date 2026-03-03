package middleware

import (
	"shrine/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func request() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := meta.BuildRequest(c)
		c.Locals(meta.RequestKey, req)
		return c.Next()
	}
}
