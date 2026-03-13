package controllers

import (
	"shrine/services"
	"shrine/utils/auth"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func StatsController(context *fiber.Ctx) error {
	auth.IsAuthenticated(context)
	citizen := auth.GetUser(context)
	return shortcuts.Success(context, services.GetStats(citizen))
}