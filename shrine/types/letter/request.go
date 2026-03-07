package letter

type CreateRequest struct {
	Recipients []string `json:"recipients"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
}

type SendMessageRequest struct {
	Body string `json:"body"`
}

type EditMessageRequest struct {
	Body string `json:"body"`
}

type RenameRequest struct {
	Title string `json:"title"`
}

type RemoveParticipantRequest struct {
	Username string `json:"username"`
}