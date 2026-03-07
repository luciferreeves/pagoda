package hypertext

import "shrine/enums"

type ServiceError struct {
	Kind    enums.ErrorKind
	Message string
}

func (e *ServiceError) Error() string { return e.Message }