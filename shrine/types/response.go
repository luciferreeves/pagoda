package types

type ErrorResponse struct {
	Error string `json:"error"`
}

type HelloResponse struct {
	Message string `json:"message"`
}
