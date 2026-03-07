package letter

type CreateRequest struct {
	Recipients     []string `json:"recipients"`
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	AttachmentRefs []string `json:"attachment_refs"`
}

type SendMessageRequest struct {
	Body           string   `json:"body"`
	AttachmentRefs []string `json:"attachment_refs"`
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