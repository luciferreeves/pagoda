package types

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	Bio         string     `json:"bio"`
	Birthday    *time.Time `json:"birthday"`
	AvatarURL   string     `json:"avatar_url"`
	BlinkieURL  string     `json:"blinkie_url"`
	Website     string     `json:"website"`
	Location    string     `json:"location"`
	Pronouns    string     `json:"pronouns"`
	Signature   string     `json:"signature"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CitizenSummary struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type StatsResponse struct {
	Citizens       int64            `json:"citizens"`
	Online         int64            `json:"online"`
	NewestCitizens []CitizenSummary `json:"newest_citizens"`
	OnlineCitizens []CitizenSummary `json:"online_citizens"`
}

type AdminUserResponse struct {
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	DisplayName     string     `json:"display_name"`
	AvatarURL       string     `json:"avatar_url"`
	Role            string     `json:"role"`
	EmailVerified   bool       `json:"email_verified"`
	AccountBanned   bool       `json:"account_banned"`
	BannedReason    string     `json:"banned_reason"`
	BannedAt        *time.Time `json:"banned_at"`
	AccountDisabled bool       `json:"account_disabled"`
	DisabledReason  string     `json:"disabled_reason"`
	DisabledAt      *time.Time `json:"disabled_at"`
	LastSeenAt      *time.Time `json:"last_seen_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PaginatedResponse struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}
