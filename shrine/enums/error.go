package enums

type ErrorKind string

const (
	BadRequest     ErrorKind = "bad_request"
	Unauthorized   ErrorKind = "unauthorized"
	Forbidden      ErrorKind = "forbidden"
	NotFound       ErrorKind = "not_found"
	Conflict       ErrorKind = "conflict"
	Unprocessable  ErrorKind = "unprocessable"
	TooManyRequest ErrorKind = "too_many_requests"
	Internal       ErrorKind = "internal"
)