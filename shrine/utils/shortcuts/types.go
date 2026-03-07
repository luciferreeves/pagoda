package shortcuts

import "github.com/gofiber/fiber/v2"

type Response struct {
	Context *fiber.Ctx
	Data    any
	Status  int
}