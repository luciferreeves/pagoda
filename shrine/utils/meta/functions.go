package meta

import (
	"shrine/types"

	"github.com/gofiber/fiber/v2"
)

func buildQueryParams(context *fiber.Ctx) []types.HTTPParam {
	params := make([]types.HTTPParam, 0)
	context.Request().URI().QueryArgs().VisitAll(func(k, v []byte) {
		params = append(params, types.HTTPParam{
			Key:   string(k),
			Value: string(v),
		})
	})
	return params
}

func buildRouteParams(context *fiber.Ctx) []types.HTTPParam {
	params := make([]types.HTTPParam, 0)
	for k, v := range context.AllParams() {
		params = append(params, types.HTTPParam{
			Key:   k,
			Value: v,
		})
	}
	return params
}

func buildHeaders(context *fiber.Ctx) []types.HTTPParam {
	params := make([]types.HTTPParam, 0)
	context.Request().Header.VisitAll(func(k, v []byte) {
		params = append(params, types.HTTPParam{
			Key:   string(k),
			Value: string(v),
		})
	})
	return params
}
