package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	TokenHash string `gorm:"uniqueIndex;size:64;not null"`
	UserID    uint   `gorm:"index;not null"`
	User      User   `gorm:"foreignKey:UserID"`
	ExpiresAt time.Time
	IPAddress string `gorm:"size:45"`
	UserAgent string `gorm:"size:512"`
}

func (self *Token) BeforeCreate(tx *gorm.DB) error {
	if self.TokenHash == "" {
		return gorm.ErrInvalidData
	}
	if self.UserID == 0 {
		return gorm.ErrInvalidData
	}
	if self.ExpiresAt.IsZero() {
		return gorm.ErrInvalidData
	}
	return nil
}