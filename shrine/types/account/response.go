package account

import "shrine/types/user"

type AuthResponse struct {
	Token string            `json:"token"`
	User  user.UserResponse `json:"user"`
}