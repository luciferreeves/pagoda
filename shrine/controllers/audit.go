package controllers

import (
	"shrine/services"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ListAuditLogsController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	action, _ := meta.Request(context).Query("action")
	targetType, _ := meta.Request(context).Query("target_type")

	items, total := services.ListAuditLogs(pagination, action, targetType)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func GetAuditLogController(context *fiber.Ctx) error {
	ref := meta.Request(context).MustHave().Param("ref")

	result, serviceErr := services.GetAuditLog(ref)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}