package meta

import (
	"shrine/types"
	"shrine/utils/logger"

	"github.com/gofiber/fiber/v2"
)

const RequestKey = "__request_context"

func Request(context *fiber.Ctx) facade {
	request, ok := context.Locals(RequestKey).(types.Request)
	if !ok {
		logger.Errorf("META", "RequestContext missing in fiber locals")
		return facade{}
	}
	return facade{request: request, context: context}
}

func (f facade) Param(key string) (string, bool) {
	if f.context != nil {
		val := f.context.Params(key)
		if val != "" {
			return val, true
		}
	}
	return "", false
}

func (f facade) Query(key string) (string, bool) {
	for _, q := range f.request.Query {
		if q.Key == key {
			return q.Value, true
		}
	}
	return "", false
}

func (f facade) Header(key string) (string, bool) {
	for _, h := range f.request.Headers {
		if h.Key == key {
			return h.Value, true
		}
	}
	return "", false
}

func (r required) Param(key string) string {
	// Access params directly from fiber context (available after route matching)
	if r.context != nil {
		val := r.context.Params(key)
		if val != "" {
			return val
		}
	}
	logger.Errorf("META", "missing required param: %s", key)
	return ""
}

func (r required) Query(key string) string {
	for _, q := range r.request.Query {
		if q.Key == key {
			return q.Value
		}
	}
	logger.Errorf("META", "missing required query: %s", key)
	return ""
}

func (r required) Header(key string) string {
	for _, h := range r.request.Headers {
		if h.Key == key {
			return h.Value
		}
	}
	logger.Errorf("META", "missing required header: %s", key)
	return ""
}

func (d withDefault) Param(key string) string {
	if d.context != nil {
		val := d.context.Params(key)
		if val != "" {
			return val
		}
	}
	return d.defaults
}

func (d withDefault) Query(key string) string {
	for _, q := range d.request.Query {
		if q.Key == key {
			return q.Value
		}
	}
	return d.defaults
}

func (d withDefault) Header(key string) string {
	for _, h := range d.request.Headers {
		if h.Key == key {
			return h.Value
		}
	}
	return d.defaults
}
