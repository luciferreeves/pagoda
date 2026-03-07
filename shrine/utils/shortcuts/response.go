package shortcuts

import (
	"shrine/enums"
	"shrine/types/common"
	"shrine/types/hypertext"

	"github.com/gofiber/fiber/v2"
)

func Respond(context *fiber.Ctx, data any) *Response {
	return &Response{
		Context: context,
		Data:    data,
		Status:  fiber.StatusOK,
	}
}

func (self *Response) As(status int) error {
	self.Status = status
	if self.Data == nil {
		return self.Context.SendStatus(status)
	}
	return self.Context.Status(status).JSON(self.Data)
}

func HandleError(context *fiber.Ctx, err *hypertext.ServiceError) error {
	statusMap := map[enums.ErrorKind]int{
		enums.BadRequest:     fiber.StatusBadRequest,
		enums.Unauthorized:   fiber.StatusUnauthorized,
		enums.Forbidden:      fiber.StatusForbidden,
		enums.NotFound:       fiber.StatusNotFound,
		enums.Conflict:       fiber.StatusConflict,
		enums.Unprocessable:  fiber.StatusUnprocessableEntity,
		enums.TooManyRequest: fiber.StatusTooManyRequests,
		enums.Internal:       fiber.StatusInternalServerError,
	}

	status, exists := statusMap[err.Kind]
	if !exists {
		status = fiber.StatusInternalServerError
	}

	return Respond(context, common.ErrorResponse{
		Error: err.Message,
	}).As(status)
}

func BadRequest(context *fiber.Ctx, err error) error {
	return Respond(context, common.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusBadRequest)
}

func Unauthorized(context *fiber.Ctx, err error) error {
	return Respond(context, common.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusUnauthorized)
}

func Forbidden(context *fiber.Ctx, err error) error {
	return Respond(context, common.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusForbidden)
}

func NotFound(context *fiber.Ctx, err error) error {
	return Respond(context, common.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusNotFound)
}

func InternalServerError(context *fiber.Ctx, err error) error {
	return Respond(context, common.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusInternalServerError)
}

func Success(context *fiber.Ctx, data any) error {
	return Respond(context, data).As(fiber.StatusOK)
}

func Created(context *fiber.Ctx, data any) error {
	return Respond(context, data).As(fiber.StatusCreated)
}

func NoContent(context *fiber.Ctx) error {
	return Respond(context, nil).As(fiber.StatusNoContent)
}