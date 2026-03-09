package models

import (
	"shrine/types/letter"
	"shrine/utils/crypto"
	"shrine/utils/storage"
	"time"

	"gorm.io/gorm"
)

type Letter struct {
	gorm.Model
	Ref       string `gorm:"size:12;uniqueIndex;not null"`
	Title     string `gorm:"size:200"`
	IsSystem  bool   `gorm:"not null;default:false"`
	SystemRef string `gorm:"size:20;index"`
	CreatorID *uint  `gorm:"index"`
	Creator   *User  `gorm:"foreignKey:CreatorID"`
}

type LetterParticipant struct {
	gorm.Model
	LetterID   uint       `gorm:"uniqueIndex:idx_letter_user;not null"`
	Letter     Letter     `gorm:"foreignKey:LetterID"`
	UserID     uint       `gorm:"uniqueIndex:idx_letter_user;index;not null"`
	User       User       `gorm:"foreignKey:UserID"`
	Role       string     `gorm:"size:10;not null;default:member"`
	LastReadAt *time.Time
}

type LetterMessage struct {
	gorm.Model
	Ref         string             `gorm:"size:12;uniqueIndex;not null"`
	LetterID    uint               `gorm:"index;not null"`
	Letter      Letter             `gorm:"foreignKey:LetterID"`
	SenderID    *uint              `gorm:"index"`
	Sender      *User              `gorm:"foreignKey:SenderID"`
	Body        string             `gorm:"type:text"`
	Attachments []LetterAttachment `gorm:"foreignKey:MessageID"`
	EditedAt    *time.Time
}

type LetterAttachment struct {
	gorm.Model
	Ref         string `gorm:"size:12;uniqueIndex;not null"`
	MessageID   *uint  `gorm:"index"`
	UploaderID  uint   `gorm:"index;not null"`
	FileName    string `gorm:"size:255;not null"`
	FilePath    string `gorm:"size:512;not null"`
	FileSize    int64  `gorm:"not null"`
	ContentType string `gorm:"size:100;not null"`
	Category    string `gorm:"size:20;not null;default:other"`
}

func (self *Letter) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *LetterMessage) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *LetterAttachment) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *LetterParticipant) ToResponse() letter.ParticipantResponse {
	return letter.ParticipantResponse{
		Username:    self.User.Username,
		DisplayName: self.User.DisplayName,
		AvatarURL:   storage.ResolveCDN(self.User.AvatarURL),
		Role:        self.Role,
	}
}

func (self *LetterAttachment) ToResponse() letter.AttachmentResponse {
	return letter.AttachmentResponse{
		Ref:         self.Ref,
		FileName:    self.FileName,
		URL:         storage.ResolveCDN(self.FilePath),
		FileSize:    self.FileSize,
		ContentType: self.ContentType,
		Category:    self.Category,
	}
}

func (self *LetterMessage) ToResponse() letter.MessageResponse {
	response := letter.MessageResponse{
		Ref:       self.Ref,
		Body:      self.Body,
		EditedAt:  self.EditedAt,
		CreatedAt: self.CreatedAt,
		Deleted:   self.DeletedAt.Valid,
	}

	if self.IsDeleted() {
		response.Body = ""
		response.Attachments = []letter.AttachmentResponse{}
	} else {
		attachments := make([]letter.AttachmentResponse, len(self.Attachments))
		for index, attach := range self.Attachments {
			attachments[index] = attach.ToResponse()
		}
		response.Attachments = attachments
	}

	if self.Sender != nil {
		summary := self.Sender.ToSummary()
		response.Sender = &summary
	}

	return response
}

func (self *LetterMessage) IsDeleted() bool {
	return self.DeletedAt.Valid
}