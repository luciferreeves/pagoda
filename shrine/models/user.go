package models

import (
	"fmt"
	"shrine/enums"
	"shrine/types"
	"shrine/utils/validators"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string         `gorm:"uniqueIndex;size:32;not null"`
	Email           string         `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash    string         `gorm:"size:255;not null"`
	DisplayName     string         `gorm:"size:50;not null"`
	Bio             string         `gorm:"size:500"`
	Birthday        *time.Time
	AvatarURL       string         `gorm:"size:512"`
	BlinkieURL      string         `gorm:"size:512"`
	Website         string         `gorm:"size:255"`
	Location        string         `gorm:"size:100"`
	Pronouns        string         `gorm:"size:50"`
	Signature       string         `gorm:"size:500"`
	Role            enums.UserRole `gorm:"size:20;not null;default:member"`
	EmailVerified      bool       `gorm:"not null;default:false"`
	VerificationHash   string     `gorm:"size:64"`
	VerificationExpiry *time.Time
	VerificationType   enums.VerificationType `gorm:"size:20"`
	AccountBanned      bool       `gorm:"not null;default:false"`
	BannedAt        *time.Time
	BannedReason    string `gorm:"size:500"`
	BannedBy        *uint  `gorm:"index"`
	AccountDisabled bool   `gorm:"not null;default:false"`
	DisabledAt      *time.Time
	DisabledReason  string `gorm:"size:500"`
	DisabledBy      *uint  `gorm:"index"`
	LastSeenAt      *time.Time
	RegistrationIP  string `gorm:"size:45"`
}

func (user *User) SetPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 255 {
		return fmt.Errorf("password must be at most 255 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

func (user *User) IsAdmin() bool {
	return user.Role == enums.Admin
}

func (user *User) IsModerator() bool {
	return user.Role == enums.Moderator
}

func (user *User) IsStaff() bool {
	return user.IsAdmin() || user.IsModerator()
}

func (user *User) IsBanned() bool {
	return user.AccountBanned
}

func (user *User) IsDisabled() bool {
	return user.AccountDisabled
}

func (user *User) CanAuthenticate() bool {
	return !user.AccountBanned && !user.AccountDisabled
}

func (user *User) IsVerified() bool {
	return user.EmailVerified
}

func (user *User) SetVerification(hash string, expiry time.Time, verificationType enums.VerificationType) {
	user.VerificationHash = hash
	user.VerificationExpiry = &expiry
	user.VerificationType = verificationType
}

func (user *User) VerifyEmail() {
	user.EmailVerified = true
	user.VerificationHash = ""
	user.VerificationExpiry = nil
	user.VerificationType = ""
}

func (user *User) ToResponse() types.UserResponse {
	return types.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		Birthday:    user.Birthday,
		AvatarURL:   user.AvatarURL,
		BlinkieURL:  user.BlinkieURL,
		Website:     user.Website,
		Location:    user.Location,
		Pronouns:    user.Pronouns,
		Signature:   user.Signature,
		Role:        string(user.Role),
		CreatedAt:   user.CreatedAt,
	}
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	_, bypassUsername := tx.Get("bypass_username_validation")

	if !bypassUsername {
		if !validators.IsValidUsername(user.Username, 3) {
			return fmt.Errorf("username must be 3-32 characters, alphanumeric and underscores only")
		}
		if validators.IsReservedUsername(user.Username) {
			return fmt.Errorf("username is reserved")
		}
	}

	if !validators.IsValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	if len(strings.TrimSpace(user.DisplayName)) < 1 || len(strings.TrimSpace(user.DisplayName)) > 50 {
		return fmt.Errorf("display name must be between 1 and 50 characters")
	}

	if user.PasswordHash == "" {
		return fmt.Errorf("password is required")
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.DisplayName = strings.TrimSpace(user.DisplayName)
	user.Username = strings.TrimSpace(user.Username)

	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if !validators.IsValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	if len(strings.TrimSpace(user.DisplayName)) < 1 || len(strings.TrimSpace(user.DisplayName)) > 50 {
		return fmt.Errorf("display name must be between 1 and 50 characters")
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.DisplayName = strings.TrimSpace(user.DisplayName)

	return nil
}