package controllers

import (
	"shrine/services"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func StatsController(context *fiber.Ctx) error {
	return shortcuts.Success(context, services.GetStats())
}