package ticket

type CreateRequest struct {
	CategoryRef string `json:"category_ref"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Priority    string `json:"priority"`
}

type UpdateRequest struct {
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	CategoryRef *string `json:"category_ref"`
	Assignee    *string `json:"assignee"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	SortOrder   *uint   `json:"sort_order"`
}

type SendMessageRequest struct {
	Body string `json:"body"`
}