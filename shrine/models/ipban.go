package models

import (
	"time"
)

type IPBan struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	IP        string    `gorm:"size:45;uniqueIndex;not null"`
	Reason    string    `gorm:"size:500"`
	CreatedAt time.Time `gorm:"index"`
}