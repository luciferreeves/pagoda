package meta

import (
	"shrine/types"

	"github.com/gofiber/fiber/v2"
)

func BuildRequest(context *fiber.Ctx) types.Request {
	return types.Request{
		Path:        context.Path(),
		Method:      context.Method(),
		Query:       buildQueryParams(context),
		Params:      buildRouteParams(context),
		Headers:     buildHeaders(context),
		QueryString: string(context.Request().URI().QueryString()),
		IP:          context.IP(),
		URL:         context.OriginalURL(),
	}
}
