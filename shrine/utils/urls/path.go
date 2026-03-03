package urls

import (
	"shrine/types"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Path(method types.HTTPMethod, path string, handler fiber.Handler, name string) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	namespace := registry.currentNamespace
	fullName := name
	fullPath := path

	if namespace != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		fullName = namespace + "." + name
		fullPath = "/" + namespace + path
	} else {
		if !strings.HasPrefix(fullPath, "/") {
			fullPath = "/" + fullPath
		}
	}

	registry.routes[fullName] = registeredRoute{
		method:    method,
		path:      path,
		handler:   handler,
		namespace: namespace,
		name:      name,
		fullPath:  fullPath,
	}
}

func GetFullPath(routeName string) (string, bool) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	route, ok := registry.routes[routeName]
	if !ok {
		return "", false
	}

	return route.fullPath, true
}
