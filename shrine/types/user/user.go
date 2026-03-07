package user

import "time"

type CitizenSummaryResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
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

type AdminUserResponse struct {
	UserResponse
	Jade            uint64     `json:"jade"`
	Honor           uint64     `json:"honor"`
	EmailVerified   bool       `json:"email_verified"`
	WarningCount    uint       `json:"warning_count"`
	AccountBanned   bool       `json:"account_banned"`
	BannedReason    string     `json:"banned_reason"`
	BannedAt        *time.Time `json:"banned_at"`
	AccountDisabled bool       `json:"account_disabled"`
	DisabledReason  string     `json:"disabled_reason"`
	DisabledAt      *time.Time `json:"disabled_at"`
	DisabledUntil   *time.Time `json:"disabled_until"`
	LastSeenAt      *time.Time `json:"last_seen_at"`
	RegistrationIP  string     `json:"registration_ip"`
}

type StatsResponse struct {
	Citizens       int64                    `json:"citizens"`
	Online         int64                    `json:"online"`
	NewestCitizens []CitizenSummaryResponse `json:"newest_citizens"`
	OnlineCitizens []CitizenSummaryResponse `json:"online_citizens"`
}