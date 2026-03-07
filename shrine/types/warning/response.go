package warning

import "time"

type WarningResponse struct {
	SystemRef string    `json:"system_ref"`
	Admin     string    `json:"admin"`
	Message   string    `json:"message"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}