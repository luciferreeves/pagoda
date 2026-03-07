package urls

import (
	"shrine/enums"
	"shrine/utils/logger"

	"github.com/gofiber/fiber/v2"
)

var methodBinders = map[enums.HTTPMethod]func(fiber.Router, string, fiber.Handler) fiber.Router{
	enums.GET:     func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Get(path, h) },
	enums.POST:    func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Post(path, h) },
	enums.PUT:     func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Put(path, h) },
	enums.PATCH:   func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Patch(path, h) },
	enums.DELETE:  func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Delete(path, h) },
	enums.OPTIONS: func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Options(path, h) },
	enums.HEAD:    func(r fiber.Router, path string, h fiber.Handler) fiber.Router { return r.Head(path, h) },
}

func Attach(app *fiber.App) {
	namespaceGroups := make(map[string]fiber.Router)

	for fullName, route := range registry.routes {
		group, exists := namespaceGroups[route.namespace]
		if !exists {
			group = app.Group("/" + route.namespace)
			namespaceGroups[route.namespace] = group
		}

		binder, ok := methodBinders[route.method]
		if !ok {
			logger.Fatalf("URLs", "Unsupported HTTP method: %s for route %s", route.method, fullName)
		}

		fiberRoute := binder(group, route.path, route.handler)
		fiberRoute.Name(fullName)
	}
}
