package models

import (
	"errors"
	"shrine/enums"
	"shrine/messages"
	"shrine/types/user"
	"shrine/utils/storage"
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
	DisplayName string         `gorm:"size:50;not null"`
	Bio         string         `gorm:"size:500"`
	Birthday    *time.Time
	AvatarURL   string         `gorm:"size:512;not null;default:defaults/avatar.png"`
	BlinkieURL  string         `gorm:"size:512"`
	Website     string         `gorm:"size:255"`
	Location    string         `gorm:"size:100"`
	Pronouns    string         `gorm:"size:50"`
	Signature   string         `gorm:"size:500"`
	Jade        uint64         `gorm:"not null;default:0"`
	Honor       uint64         `gorm:"not null;default:0"`
	Role            enums.UserRole         `gorm:"size:20;not null;default:member"`
	EmailVerified      bool                `gorm:"not null;default:false"`
	VerificationHash   string              `gorm:"size:64"`
	VerificationExpiry *time.Time
	VerificationType   enums.VerificationType `gorm:"size:20"`
	AccountBanned   bool       `gorm:"not null;default:false"`
	BannedAt        *time.Time
	BannedReason    string     `gorm:"size:500"`
	BannedBy        *uint      `gorm:"index"`
	AccountDisabled bool       `gorm:"not null;default:false"`
	DisabledAt      *time.Time
	DisabledReason  string     `gorm:"size:500"`
	DisabledBy      *uint      `gorm:"index"`
	DisabledUntil   *time.Time
	WarningCount    uint       `gorm:"not null;default:0"`
	LastSeenAt      *time.Time
	IP              string     `gorm:"size:45"`
}

func (self *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New(messages.PasswordTooShort)
	}
	if len(password) > 255 {
		return errors.New(messages.PasswordTooLong)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	self.PasswordHash = string(hash)
	return nil
}

func (self *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(self.PasswordHash), []byte(password)) == nil
}

func (self *User) IsOwner() bool {
	return self.Role == enums.Owner
}

func (self *User) IsAdmin() bool {
	return self.Role == enums.Admin || self.IsOwner()
}

func (self *User) IsModerator() bool {
	return self.Role == enums.Moderator
}

func (self *User) IsStaff() bool {
	return self.IsAdmin() || self.IsModerator()
}

func (self *User) CanAuthenticate() bool {
	return !self.AccountBanned && !self.AccountDisabled
}

func (self *User) ClearExpiredDisable() bool {
	if !self.AccountDisabled || self.DisabledUntil == nil || !self.DisabledUntil.Before(time.Now()) {
		return false
	}

	self.AccountDisabled = false
	self.DisabledAt = nil
	self.DisabledReason = ""
	self.DisabledBy = nil
	self.DisabledUntil = nil
	return true
}

func (self *User) IsVerified() bool {
	return self.EmailVerified
}

func (self *User) SetVerification(hash string, expiry time.Time, verificationType enums.VerificationType) {
	self.VerificationHash = hash
	self.VerificationExpiry = &expiry
	self.VerificationType = verificationType
}

func (self *User) VerifyEmail() {
	self.EmailVerified = true
	self.VerificationHash = ""
	self.VerificationExpiry = nil
	self.VerificationType = ""
}

func (self *User) ToResponse() user.UserResponse {
	return user.UserResponse{
		Username:    self.Username,
		Email:       self.Email,
		DisplayName: self.DisplayName,
		Bio:         self.Bio,
		Birthday:    self.Birthday,
		AvatarURL:   storage.ResolveCDN(self.AvatarURL),
		BlinkieURL:  storage.ResolveCDN(self.BlinkieURL),
		Website:     self.Website,
		Location:    self.Location,
		Pronouns:    self.Pronouns,
		Signature:   self.Signature,
		Role:        string(self.Role),
		CreatedAt:   self.CreatedAt,
	}
}

func (self *User) ToAdminResponse() user.AdminUserResponse {
	return user.AdminUserResponse{
		UserResponse:    self.ToResponse(),
		Jade:            self.Jade,
		Honor:           self.Honor,
		EmailVerified:   self.EmailVerified,
		WarningCount:    self.WarningCount,
		AccountBanned:   self.AccountBanned,
		BannedReason:    self.BannedReason,
		BannedAt:        self.BannedAt,
		AccountDisabled: self.AccountDisabled,
		DisabledReason:  self.DisabledReason,
		DisabledAt:      self.DisabledAt,
		DisabledUntil:   self.DisabledUntil,
		LastSeenAt:      self.LastSeenAt,
	}
}

func (self *User) ToSummary() user.CitizenSummaryResponse {
	return user.CitizenSummaryResponse{
		Username:    self.Username,
		DisplayName: self.DisplayName,
		AvatarURL:   storage.ResolveCDN(self.AvatarURL),
	}
}

func (self *User) BeforeCreate(tx *gorm.DB) error {
	_, bypassUsername := tx.Get("bypass_username_validation")

	if !bypassUsername {
		if !validators.IsValidUsername(self.Username, 3) {
			return errors.New(messages.InvalidUsername)
		}
		if validators.IsReservedUsername(self.Username) {
			return errors.New(messages.UsernameNotAvailable)
		}
	}

	if !validators.IsValidEmail(self.Email) {
		return errors.New(messages.InvalidEmail)
	}

	if len(strings.TrimSpace(self.DisplayName)) < 1 || len(strings.TrimSpace(self.DisplayName)) > 50 {
		return errors.New(messages.InvalidDisplayName)
	}

	if self.PasswordHash == "" {
		return errors.New(messages.PasswordRequired)
	}

	self.Email = strings.ToLower(strings.TrimSpace(self.Email))
	self.DisplayName = strings.TrimSpace(self.DisplayName)
	self.Username = strings.TrimSpace(self.Username)

	return nil
}

func (self *User) BeforeUpdate(tx *gorm.DB) error {
	if !validators.IsValidEmail(self.Email) {
		return errors.New(messages.InvalidEmail)
	}

	if len(strings.TrimSpace(self.DisplayName)) < 1 || len(strings.TrimSpace(self.DisplayName)) > 50 {
		return errors.New(messages.InvalidDisplayName)
	}

	if self.Jade > validators.MaxJade {
		return errors.New(messages.JadeExceedsMax)
	}

	self.Email = strings.ToLower(strings.TrimSpace(self.Email))
	self.DisplayName = strings.TrimSpace(self.DisplayName)

	return nil
}