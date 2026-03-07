package warning

type WarnRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}