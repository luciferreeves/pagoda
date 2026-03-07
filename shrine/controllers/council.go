package controllers

import (
	"shrine/services"
	"shrine/types/council"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

var userSortColumns = []string{"username", "display_name", "email", "role", "created_at"}

func ListUsersController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	sorting := meta.Sort(context, userSortColumns, "created_at")
	search, _ := meta.Request(context).Query("search")

	items, total := services.ListUsers(pagination, sorting, search)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func GetUserController(context *fiber.Ctx) error {
	username := meta.Request(context).MustHave().Param("username")

	result, serviceErr := services.GetUser(username)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func BanUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	body, _ := meta.Body[council.BanRequest](context)

	result, serviceErr := services.BanUser(admin, target, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func UnbanUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	result, serviceErr := services.UnbanUser(admin, target)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func DisableUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	body, _ := meta.Body[council.DisableRequest](context)

	result, serviceErr := services.DisableUser(admin, target, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func EnableUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	result, serviceErr := services.EnableUser(admin, target)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func ChangeRoleController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	body, err := meta.Body[council.ChangeRoleRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.ChangeRole(admin, target, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func EditUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	target, serviceErr := services.ResolveUser(meta.Request(context).MustHave().Param("username"))
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	body, err := meta.Body[council.EditUserRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.EditUser(admin, target, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}