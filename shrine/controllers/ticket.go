package controllers

import (
	"shrine/services"
	"shrine/types/ticket"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ListUserTicketsController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	pagination := meta.Paginate(context)

	items, total := services.ListUserTickets(citizen.ID, pagination)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func GetUserTicketController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	result, serviceErr := services.GetUserTicket(ref, citizen.ID)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func CreateTicketController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)

	body, err := meta.Body[ticket.CreateRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.CreateTicket(citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func ReplyTicketController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[ticket.SendMessageRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.ReplyTicket(ref, citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func ListAllTicketsController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	status, _ := meta.Request(context).Query("status")
	priority, _ := meta.Request(context).Query("priority")

	items, total := services.ListAllTickets(pagination, status, priority)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func GetStaffTicketController(context *fiber.Ctx) error {
	ref := meta.Request(context).MustHave().Param("ref")

	result, serviceErr := services.GetStaffTicket(ref)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func UpdateTicketController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[ticket.UpdateRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.UpdateTicket(admin.ID, ref, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func StaffReplyTicketController(context *fiber.Ctx) error {
	staff := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[ticket.SendMessageRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.StaffReplyTicket(ref, staff.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func ListTicketCategoriesController(context *fiber.Ctx) error {
	return shortcuts.Success(context, services.ListTicketCategories())
}

func CreateTicketCategoryController(context *fiber.Ctx) error {
	body, err := meta.Body[ticket.CreateCategoryRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.CreateTicketCategory(body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func UpdateTicketCategoryController(context *fiber.Ctx) error {
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[ticket.UpdateCategoryRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.UpdateTicketCategory(ref, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func DeleteTicketCategoryController(context *fiber.Ctx) error {
	ref := meta.Request(context).MustHave().Param("ref")

	serviceErr := services.DeleteTicketCategory(ref)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.NoContent(context)
}