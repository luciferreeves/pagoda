package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"shrine/config"
	"shrine/messages"
	"shrine/repositories"
	"shrine/services"
	"shrine/types/council"
	"shrine/utils/auth"
	"shrine/utils/crypto"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"
	"shrine/utils/storage"

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

func ListIPBansController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	sorting := meta.Sort(context, []string{"ip", "reason", "created_at"}, "created_at")
	items, total := repositories.ListIPBans(pagination, sorting)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func DeleteIPBanController(context *fiber.Ctx) error {
	id, err := strconv.ParseUint(meta.Request(context).MustHave().Param("id"), 10, 64)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}
	repositories.DeleteIPBanByID(uint(id))
	return shortcuts.NoContent(context)
}

func UploadImageController(context *fiber.Ctx) error {
	file, err := context.FormFile("file")
	if err != nil {
		return shortcuts.BadRequest(context, fiber.NewError(fiber.StatusBadRequest, messages.FileRequired))
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return shortcuts.BadRequest(context, fiber.NewError(fiber.StatusBadRequest, messages.OnlyImagesAllowed))
	}

	if file.Size > int64(config.Storage.MaxFileSize) {
		return shortcuts.BadRequest(context, fiber.NewError(fiber.StatusBadRequest, messages.FileTooLarge))
	}

	source, err := file.Open()
	if err != nil {
		return shortcuts.InternalServerError(context, err)
	}
	defer source.Close()

	ref := crypto.Ref()
	path := fmt.Sprintf("citizens/signatures/%s/%s", ref, file.Filename)

	if err := storage.Upload(path, source, file.Size, contentType); err != nil {
		return shortcuts.InternalServerError(context, err)
	}

	return shortcuts.Created(context, fiber.Map{"url": storage.ResolveCDN(path)})
}