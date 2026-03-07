package services

import (
	"shrine/enums"
	"shrine/types/hypertext"
)

func fail(kind enums.ErrorKind, message string) *hypertext.ServiceError {
	return &hypertext.ServiceError{Kind: kind, Message: message}
}