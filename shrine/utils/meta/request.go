package meta

import (
	"shrine/types"
	"shrine/utils/logger"

	"github.com/gofiber/fiber/v2"
)

const RequestKey = "__request_ctx"

func Request(c *fiber.Ctx) facade {
	req, ok := c.Locals(RequestKey).(types.Request)
	if !ok {
		logger.Errorf("META", "RequestContext missing in fiber locals")
		return facade{}
	}
	return facade{req: req, ctx: c}
}

func (f facade) Param(key string) (string, bool) {
	if f.ctx != nil {
		val := f.ctx.Params(key)
		if val != "" {
			return val, true
		}
	}
	return "", false
}

func (f facade) Query(key string) (string, bool) {
	for _, q := range f.req.Query {
		if q.Key == key {
			return q.Value, true
		}
	}
	return "", false
}

func (f facade) Header(key string) (string, bool) {
	for _, h := range f.req.Headers {
		if h.Key == key {
			return h.Value, true
		}
	}
	return "", false
}

func (r required) Param(key string) string {
	// Access params directly from fiber context (available after route matching)
	if r.ctx != nil {
		val := r.ctx.Params(key)
		if val != "" {
			return val
		}
	}
	logger.Errorf("META", "missing required param: %s", key)
	return ""
}

func (r required) Query(key string) string {
	for _, q := range r.req.Query {
		if q.Key == key {
			return q.Value
		}
	}
	logger.Errorf("META", "missing required query: %s", key)
	return ""
}

func (r required) Header(key string) string {
	for _, h := range r.req.Headers {
		if h.Key == key {
			return h.Value
		}
	}
	logger.Errorf("META", "missing required header: %s", key)
	return ""
}

func (d withDefault) Param(key string) string {
	if d.ctx != nil {
		val := d.ctx.Params(key)
		if val != "" {
			return val
		}
	}
	return d.def
}

func (d withDefault) Query(key string) string {
	for _, q := range d.req.Query {
		if q.Key == key {
			return q.Value
		}
	}
	return d.def
}

func (d withDefault) Header(key string) string {
	for _, h := range d.req.Headers {
		if h.Key == key {
			return h.Value
		}
	}
	return d.def
}
