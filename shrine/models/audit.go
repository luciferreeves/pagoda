package models

import (
	"shrine/types/audit"
	"time"
)

type AuditLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	SystemRef  string    `gorm:"size:20;uniqueIndex"`
	ActorID    uint      `gorm:"index;not null"`
	Actor      User      `gorm:"foreignKey:ActorID"`
	Action     string    `gorm:"size:50;index;not null"`
	TargetType string    `gorm:"size:30"`
	TargetRef  string    `gorm:"size:100"`
	Summary    string    `gorm:"size:500"`
	Details    string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"index"`
}

func (self *AuditLog) ToResponse() audit.AuditLogResponse {
	return audit.AuditLogResponse{
		SystemRef:  self.SystemRef,
		Actor:      self.Actor.Username,
		Action:     self.Action,
		TargetType: self.TargetType,
		TargetRef:  self.TargetRef,
		Summary:    self.Summary,
		CreatedAt:  self.CreatedAt,
	}
}

func (self *AuditLog) ToDetailResponse() audit.DetailResponse {
	return audit.DetailResponse{
		AuditLogResponse: self.ToResponse(),
		Details:          self.Details,
	}
}