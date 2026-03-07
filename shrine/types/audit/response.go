package audit

import "time"

type AuditLogResponse struct {
	SystemRef  string    `json:"system_ref"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type"`
	TargetRef  string    `json:"target_ref"`
	Summary    string    `json:"summary"`
	CreatedAt  time.Time `json:"created_at"`
}

type DetailResponse struct {
	AuditLogResponse
	Details string `json:"details"`
}