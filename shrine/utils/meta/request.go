package meta

import (
	"shrine/types/hypertext"
	"shrine/utils/logger"

	"github.com/gofiber/fiber/v2"
)

const RequestKey = "__request_context"

func Request(context *fiber.Ctx) facade {
	request, ok := context.Locals(RequestKey).(hypertext.Request)
	if !ok {
		logger.Errorf("META", "RequestContext missing in fiber locals")
		return facade{}
	}
	return facade{Request: request, context: context}
}

func (self facade) Param(key string) (string, bool) {
	if self.context != nil {
		val := self.context.Params(key)
		if val != "" {
			return val, true
		}
	}
	return "", false
}

func (self facade) Query(key string) (string, bool) {
	return findParam(self.Request.Query, key)
}

func (self facade) Header(key string) (string, bool) {
	return findParam(self.Request.Headers, key)
}

func (self required) Param(key string) string {
	if self.context != nil {
		val := self.context.Params(key)
		if val != "" {
			return val
		}
	}
	logger.Errorf("META", "missing required param: %s", key)
	return ""
}

func (self required) Query(key string) string {
	value, found := findParam(self.request.Query, key)
	if !found {
		logger.Errorf("META", "missing required query: %s", key)
	}
	return value
}

func (self required) Header(key string) string {
	value, found := findParam(self.request.Headers, key)
	if !found {
		logger.Errorf("META", "missing required header: %s", key)
	}
	return value
}

func (self withDefault) Param(key string) string {
	if self.context != nil {
		val := self.context.Params(key)
		if val != "" {
			return val
		}
	}
	return self.defaults
}

func (self withDefault) Query(key string) string {
	value, found := findParam(self.request.Query, key)
	if found {
		return value
	}
	return self.defaults
}

func (self withDefault) Header(key string) string {
	value, found := findParam(self.request.Headers, key)
	if found {
		return value
	}
	return self.defaults
}