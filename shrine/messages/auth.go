package messages

const (
	InvalidRequestBody          = "Invalid request body."
	InvalidUsernameOrPassword   = "Invalid username or password."
	AccountBannedOrDisabled     = "Your account has been banned or disabled."
	EmailNotVerified            = "Your email address has not been verified. Please check your inbox."
	VerificationTokenRequired   = "Verification token is required."
	VerificationLinkInvalid     = "Your verification link is invalid or has expired."
	InvalidVerificationType     = "Invalid verification type."
	AccountAlreadyVerified      = "This account has already been verified."
	NoAccountWithEmail          = "No account exists with that email address."
	UsernameAlreadyExists       = "An account with that username already exists."
	EmailAlreadyExists          = "An account with that email address already exists."
	AccountCreated              = "Your account has been created. Please check your email to verify your account."
	EmailVerified               = "Your email has been verified successfully. You can now log in."
	VerificationEmailSent       = "A new verification email has been sent. Please check your inbox."
	LoggedOut                   = "You have been logged out successfully."
	FailedGenerateToken         = "Failed to generate verification token."
	FailedStoreToken            = "Failed to store verification token."
	FailedCreateSession         = "Failed to create session."
	FailedVerifyAccount         = "Failed to verify your account."
	FailedEndSession            = "Failed to end your session."
)