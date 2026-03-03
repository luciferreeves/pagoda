package types

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID          uint       `json:"id"`
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
