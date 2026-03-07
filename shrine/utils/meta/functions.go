package meta

import (
	"shrine/types/hypertext"

	"github.com/gofiber/fiber/v2"
)

func buildQueryParams(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	context.Request().URI().QueryArgs().VisitAll(func(k, v []byte) {
		params = append(params, hypertext.Param{
			Key:   string(k),
			Value: string(v),
		})
	})
	return params
}

func buildRouteParams(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	for k, v := range context.AllParams() {
		params = append(params, hypertext.Param{
			Key:   k,
			Value: v,
		})
	}
	return params
}

func buildHeaders(context *fiber.Ctx) []hypertext.Param {
	params := make([]hypertext.Param, 0)
	context.Request().Header.VisitAll(func(k, v []byte) {
		params = append(params, hypertext.Param{
			Key:   string(k),
			Value: string(v),
		})
	})
	return params
}
