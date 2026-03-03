package controllers

import (
	"fmt"
	"shrine/types"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func HelloController(context *fiber.Ctx) error {
	name := meta.Request(context).Default("World").Query("name")
	return shortcuts.Response(context, types.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
	}).As(fiber.StatusOK)
}
