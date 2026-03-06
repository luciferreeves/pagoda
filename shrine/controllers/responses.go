package controllers

import (
	"shrine/types"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusBadRequest)
}

func Unauthorized(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusUnauthorized)
}

func Forbidden(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusForbidden)
}

func NotFound(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusNotFound)
}

func InternalServerError(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusInternalServerError)
}

func DefaultError(context *fiber.Ctx, err error) error {
	return shortcuts.Response(context, types.ErrorResponse{
		Error: err.Error(),
	}).As(fiber.StatusInternalServerError)
}

func Success(context *fiber.Ctx, data any) error {
	return shortcuts.Response(context, data).As(fiber.StatusOK)
}

func Created(context *fiber.Ctx, data any) error {
	return shortcuts.Response(context, data).As(fiber.StatusCreated)
}

func NoContent(context *fiber.Ctx) error {
	return shortcuts.Response(context, nil).As(fiber.StatusNoContent)
}
