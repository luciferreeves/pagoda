package models

import (
	"shrine/types/ticket"
	"shrine/utils/crypto"

	"gorm.io/gorm"
)

type TicketCategory struct {
	gorm.Model
	Ref         string `gorm:"size:12;uniqueIndex;not null"`
	Name        string `gorm:"size:100;not null;uniqueIndex"`
	Description string `gorm:"size:500"`
	SortOrder   uint   `gorm:"not null;default:0"`
}

type Ticket struct {
	gorm.Model
	Ref        string         `gorm:"size:12;uniqueIndex;not null"`
	UserID     uint           `gorm:"index;not null"`
	User       User           `gorm:"foreignKey:UserID"`
	CategoryID uint           `gorm:"index;not null"`
	Category   TicketCategory `gorm:"foreignKey:CategoryID"`
	AssigneeID *uint          `gorm:"index"`
	Assignee   *User          `gorm:"foreignKey:AssigneeID"`
	Subject    string         `gorm:"size:200;not null"`
	Priority   string         `gorm:"size:10;not null;default:low"`
	Status     string         `gorm:"size:20;not null;default:open"`
}

type TicketMessage struct {
	gorm.Model
	Ref      string `gorm:"size:12;uniqueIndex;not null"`
	TicketID uint   `gorm:"index;not null"`
	Ticket   Ticket `gorm:"foreignKey:TicketID"`
	SenderID uint   `gorm:"not null"`
	Sender   User   `gorm:"foreignKey:SenderID"`
	Body     string `gorm:"type:text;not null"`
	IsStaff  bool   `gorm:"not null;default:false"`
}

func (self *TicketCategory) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *Ticket) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *TicketMessage) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *TicketCategory) ToResponse() ticket.CategoryResponse {
	return ticket.CategoryResponse{
		Ref:         self.Ref,
		Name:        self.Name,
		Description: self.Description,
		SortOrder:   self.SortOrder,
	}
}

func (self *Ticket) ToResponse() ticket.TicketResponse {
	response := ticket.TicketResponse{
		Ref:       self.Ref,
		Subject:   self.Subject,
		Category:  self.Category.ToResponse(),
		Priority:  self.Priority,
		Status:    self.Status,
		User:      self.User.ToSummary(),
		CreatedAt: self.CreatedAt,
		UpdatedAt: self.UpdatedAt,
	}

	if self.Assignee != nil {
		summary := self.Assignee.ToSummary()
		response.Assignee = &summary
	}

	return response
}

func (self *TicketMessage) ToResponse() ticket.MessageResponse {
	return ticket.MessageResponse{
		Ref:       self.Ref,
		Sender:    self.Sender.ToSummary(),
		Body:      self.Body,
		IsStaff:   self.IsStaff,
		CreatedAt: self.CreatedAt,
	}
}