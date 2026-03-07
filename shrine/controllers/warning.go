package controllers

import (
	"shrine/services"
	"shrine/types/warning"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func WarnUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	body, err := meta.Body[warning.WarnRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.WarnUser(admin, target, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func DeactivateWarningController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	ref := context.Params("ref")

	result, serviceErr := services.DeactivateWarning(admin, ref)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func ListWarningsController(context *fiber.Ctx) error {
	username := meta.Request(context).MustHave().Param("username")
	pagination := meta.Paginate(context)

	items, total, serviceErr := services.ListWarnings(username, pagination)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, pagination.Response(items, total))
}