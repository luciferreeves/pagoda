package audit

type BanDetails struct {
	Reason    string `json:"reason"`
	SystemRef string `json:"system_ref"`
}

type DisableDetails struct {
	Reason        string  `json:"reason"`
	DisabledUntil *string `json:"disabled_until"`
	SystemRef     string  `json:"system_ref"`
}

type RoleChangeDetails struct {
	OldRole string `json:"old_role"`
	NewRole string `json:"new_role"`
}

type WarningDetails struct {
	WarningRef string `json:"warning_ref"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

type DeactivateWarningDetails struct {
	WarningRef string `json:"warning_ref"`
}

type FieldChange struct {
	Field string `json:"field"`
	Old   any    `json:"old"`
	New   any    `json:"new"`
}

type EditUserDetails struct {
	Changes []FieldChange `json:"changes"`
}