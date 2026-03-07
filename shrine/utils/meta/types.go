package meta

import (
	"shrine/types/hypertext"

	"github.com/gofiber/fiber/v2"
)

type facade struct {
	hypertext.Request
	context *fiber.Ctx
}

type required struct {
	request hypertext.Request
	context *fiber.Ctx
}

type withDefault struct {
	request  hypertext.Request
	context  *fiber.Ctx
	defaults string
}
