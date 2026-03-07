package controllers

import (
	"shrine/messages"
	"shrine/services"
	"shrine/types/letter"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ListLettersController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	pagination := meta.Paginate(context)

	items, total := services.ListLetters(citizen.ID, pagination)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func GetLetterController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")
	pagination := meta.Paginate(context)

	result, serviceErr := services.GetLetter(ref, citizen.ID, pagination)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func CreateLetterController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)

	body, err := meta.Body[letter.CreateRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.CreateLetter(citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func SendMessageController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[letter.SendMessageRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.SendLetterMessage(ref, citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func EditMessageController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")
	messageRef := meta.Request(context).MustHave().Param("messageRef")

	body, err := meta.Body[letter.EditMessageRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.EditLetterMessage(ref, messageRef, citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func DeleteMessageController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")
	messageRef := meta.Request(context).MustHave().Param("messageRef")

	serviceErr := services.DeleteLetterMessage(ref, messageRef, citizen.ID)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.NoContent(context)
}

func RenameLetterController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[letter.RenameRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.RenameLetter(ref, citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func LeaveLetterController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	result, serviceErr := services.LeaveLetter(ref, citizen.ID)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func RemoveParticipantController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[letter.RemoveParticipantRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.RemoveLetterParticipant(ref, citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func UploadAttachmentController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)

	file, err := context.FormFile("file")
	if err != nil {
		return shortcuts.BadRequest(context, fiber.NewError(fiber.StatusBadRequest, messages.FileRequired))
	}

	source, err := file.Open()
	if err != nil {
		return shortcuts.InternalServerError(context, err)
	}
	defer source.Close()

	result, serviceErr := services.UploadLetterAttachment(citizen.ID, file.Filename, file.Size, file.Header.Get("Content-Type"), source)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}