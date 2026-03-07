package meta

import (
	"shrine/types/hypertext"

	"github.com/gofiber/fiber/v2"
)

func findParam(params []hypertext.Param, key string) (string, bool) {
	for _, param := range params {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

func buildQueryParams(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	context.Request().URI().QueryArgs().VisitAll(func(name, value []byte) {
		params = append(params, hypertext.Param{
			Key:   string(name),
			Value: string(value),
		})
	})
	return params
}

func buildRouteParams(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	for name, value := range context.AllParams() {
		params = append(params, hypertext.Param{
			Key:   name,
			Value: value,
		})
	}
	return params
}

func buildHeaders(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	context.Request().Header.VisitAll(func(name, value []byte) {
		params = append(params, hypertext.Param{
			Key:   string(name),
			Value: string(value),
		})
	})
	return params
}