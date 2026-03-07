package hypertext

import "shrine/enums"

type ServiceError struct {
	Kind    enums.ErrorKind
	Message string
}

func (self *ServiceError) Error() string { return self.Message }