package meta

import (
	"shrine/types"

	"github.com/gofiber/fiber/v2"
)

type facade struct {
	request types.Request
	context *fiber.Ctx
}

type required struct {
	request types.Request
	context *fiber.Ctx
}

type withDefault struct {
	request  types.Request
	context  *fiber.Ctx
	defaults string
}
