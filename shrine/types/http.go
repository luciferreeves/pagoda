package types

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	OPTIONS HTTPMethod = "OPTIONS"
	HEAD    HTTPMethod = "HEAD"
)

type HTTPParam struct {
	Key   string
	Value string
}

type Request struct {
	Path        string
	Method      string
	Query       []HTTPParam
	Params      []HTTPParam
	Headers     []HTTPParam
	QueryString string
	IP          string
	URL         string
}
