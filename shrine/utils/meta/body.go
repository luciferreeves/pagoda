package meta

import "github.com/gofiber/fiber/v2"

func Body[T any](context *fiber.Ctx) (T, error) {
	var body T
	if err := context.BodyParser(&body); err != nil {
		return body, err
	}
	return body, nil
}