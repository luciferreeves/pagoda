package council

type BanRequest struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type DisableRequest struct {
	Reason        string  `json:"reason"`
	Message       string  `json:"message"`
	DisabledUntil *string `json:"disabled_until"`
}

type EditUserRequest struct {
	Username    *string `json:"username"`
	DisplayName *string `json:"display_name"`
	Email       *string `json:"email"`
	Bio         *string `json:"bio"`
	Birthday    *string `json:"birthday"`
	AvatarURL   *string `json:"avatar_url"`
	BlinkieURL  *string `json:"blinkie_url"`
	Website     *string `json:"website"`
	Location    *string `json:"location"`
	Pronouns    *string `json:"pronouns"`
	Signature   *string `json:"signature"`
	Jade        *uint64 `json:"jade"`
	Honor       *uint64 `json:"honor"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}