package types

type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type ResendActivationRequest struct {
	Email string `json:"email"`
}

type BanUserRequest struct {
	Reason string `json:"reason"`
}

type DisableUserRequest struct {
	Reason string `json:"reason"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}