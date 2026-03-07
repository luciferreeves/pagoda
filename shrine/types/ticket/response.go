package ticket

import (
	"shrine/types/user"
	"time"
)

type CategoryResponse struct {
	Ref         string `json:"ref"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   uint   `json:"sort_order"`
}

type MessageResponse struct {
	Ref       string                       `json:"ref"`
	Sender    user.CitizenSummaryResponse  `json:"sender"`
	Body      string                       `json:"body"`
	IsStaff   bool                         `json:"is_staff"`
	CreatedAt time.Time                    `json:"created_at"`
}

type TicketResponse struct {
	Ref       string                        `json:"ref"`
	Subject   string                        `json:"subject"`
	Category  CategoryResponse              `json:"category"`
	Priority  string                        `json:"priority"`
	Status    string                        `json:"status"`
	User      user.CitizenSummaryResponse   `json:"user"`
	Assignee  *user.CitizenSummaryResponse  `json:"assignee"`
	CreatedAt time.Time                     `json:"created_at"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

type DetailResponse struct {
	TicketResponse
	Messages []MessageResponse `json:"messages"`
}