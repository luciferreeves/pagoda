package letter

import (
	"shrine/types/user"
	"time"
)

type ParticipantResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
}

type AttachmentResponse struct {
	Ref         string `json:"ref"`
	FileName    string `json:"file_name"`
	URL         string `json:"url"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
}

type MessageResponse struct {
	Ref         string                       `json:"ref"`
	Sender      *user.CitizenSummaryResponse `json:"sender"`
	Body        string                       `json:"body"`
	Attachments []AttachmentResponse         `json:"attachments"`
	EditedAt    *time.Time                   `json:"edited_at"`
	CreatedAt   time.Time                    `json:"created_at"`
	Deleted     bool                         `json:"deleted"`
}

type LetterResponse struct {
	Ref          string                `json:"ref"`
	Title        string                `json:"title"`
	IsSystem     bool                  `json:"is_system"`
	SystemRef    string                `json:"system_ref,omitempty"`
	Participants []ParticipantResponse `json:"participants"`
	LastMessage  *MessageResponse      `json:"last_message,omitempty"`
	Unread       bool                  `json:"unread"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

type DetailResponse struct {
	Ref          string                `json:"ref"`
	Title        string                `json:"title"`
	IsSystem     bool                  `json:"is_system"`
	SystemRef    string                `json:"system_ref,omitempty"`
	Participants []ParticipantResponse `json:"participants"`
	Messages     []MessageResponse     `json:"messages"`
	CreatedAt    time.Time             `json:"created_at"`
}