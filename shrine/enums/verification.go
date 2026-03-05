package enums

type VerificationType string

const (
	Activation    VerificationType = "activation"
	EmailChange   VerificationType = "email_change"
	PasswordReset VerificationType = "password_reset"
)