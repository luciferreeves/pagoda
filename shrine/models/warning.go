package models

import (
	"shrine/types/warning"

	"gorm.io/gorm"
)

type Warning struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null"`
	User      User   `gorm:"foreignKey:UserID"`
	AdminID   uint   `gorm:"not null"`
	Admin     User   `gorm:"foreignKey:AdminID"`
	LetterID  uint   `gorm:"not null"`
	Letter    Letter `gorm:"foreignKey:LetterID"`
	SystemRef string `gorm:"size:20;uniqueIndex"`
	Message   string `gorm:"type:text;not null"`
	Active    bool   `gorm:"not null;default:true"`
}

func (self *Warning) ToResponse() warning.WarningResponse {
	return warning.WarningResponse{
		SystemRef: self.SystemRef,
		Admin:     self.Admin.Username,
		Message:   self.Message,
		Active:    self.Active,
		CreatedAt: self.CreatedAt,
	}
}