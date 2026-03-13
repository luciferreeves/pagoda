package repositories

import (
	"shrine/database"
	"shrine/models"
	"shrine/utils/meta"
	"time"

	"gorm.io/gorm"
)

func CreateLetter(creatorID uint, title string, recipientIDs []uint, body string, attachmentRefs []string) (*models.Letter, error) {
	var letter models.Letter
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		letter = models.Letter{
			Title:     title,
			CreatorID: &creatorID,
		}
		if err := tx.Create(&letter).Error; err != nil {
			return err
		}

		participants := make([]models.LetterParticipant, 0, len(recipientIDs)+1)
		participants = append(participants, models.LetterParticipant{
			LetterID: letter.ID,
			UserID:   creatorID,
			Role:     "owner",
		})
		for _, uid := range recipientIDs {
			participants = append(participants, models.LetterParticipant{
				LetterID: letter.ID,
				UserID:   uid,
			})
		}
		if err := tx.Create(&participants).Error; err != nil {
			return err
		}

		msg := models.LetterMessage{
			LetterID: letter.ID,
			SenderID: &creatorID,
			Body:     body,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		if len(attachmentRefs) > 0 {
			tx.Model(&models.LetterAttachment{}).Where("ref IN ? AND uploader_id = ? AND message_id IS NULL", attachmentRefs, creatorID).Update("message_id", msg.ID)
		}

		return nil
	})

	return &letter, err
}

func CreateSystemLetter(recipientID uint, title string, body string, systemRef string) (*models.Letter, error) {
	var letter models.Letter
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		letter = models.Letter{
			Title:     title,
			IsSystem:  true,
			SystemRef: systemRef,
		}
		if err := tx.Create(&letter).Error; err != nil {
			return err
		}

		participant := models.LetterParticipant{
			LetterID: letter.ID,
			UserID:   recipientID,
		}
		if err := tx.Create(&participant).Error; err != nil {
			return err
		}

		msg := models.LetterMessage{
			LetterID: letter.ID,
			Body:     body,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		return nil
	})

	return &letter, err
}

func FindLetterByRef(ref string) (*models.Letter, error) {
	var letter models.Letter
	err := database.DB.Where("ref = ?", ref).First(&letter).Error
	return &letter, err
}

func IsLetterParticipant(letterID uint, userID uint) bool {
	var count int64
	database.DB.Model(&models.LetterParticipant{}).Where("letter_id = ? AND user_id = ?", letterID, userID).Count(&count)
	return count > 0
}

func GetLetterParticipants(letterID uint) []models.LetterParticipant {
	var participants []models.LetterParticipant
	database.DB.Where("letter_id = ?", letterID).Preload("User").Find(&participants)
	return participants
}

func GetLetterMessages(letterID uint, p meta.Pagination) ([]models.LetterMessage, int64) {
	var messages []models.LetterMessage
	var total int64

	query := database.DB.Unscoped().Model(&models.LetterMessage{}).Where("letter_id = ?", letterID)
	query.Count(&total)
	p.Apply(query.Order("created_at asc")).Preload("Sender").Preload("Attachments").Find(&messages)

	return messages, total
}

func CountUnreadLetters(userID uint) int64 {
	var count int64
	database.DB.Model(&models.LetterParticipant{}).
		Joins("JOIN letters ON letters.id = letter_participants.letter_id").
		Joins("JOIN letter_messages ON letter_messages.letter_id = letters.id AND letter_messages.deleted_at IS NULL").
		Where("letter_participants.user_id = ? AND (letter_participants.last_read_at IS NULL OR letter_messages.created_at > letter_participants.last_read_at)", userID).
		Distinct("letter_participants.letter_id").
		Count(&count)
	return count
}

func ListLettersForUser(userID uint, p meta.Pagination) ([]models.Letter, int64) {
	var letters []models.Letter
	var total int64

	subQuery := database.DB.Model(&models.LetterParticipant{}).Select("letter_id").Where("user_id = ?", userID)
	query := database.DB.Model(&models.Letter{}).Where("id IN (?)", subQuery)

	query.Count(&total)
	p.Apply(query.Order("updated_at desc")).Find(&letters)

	return letters, total
}

func GetLastMessage(letterID uint) *models.LetterMessage {
	var msg models.LetterMessage
	err := database.DB.Where("letter_id = ?", letterID).Preload("Sender").Preload("Attachments").Order("created_at desc").First(&msg).Error
	if err != nil {
		return nil
	}
	return &msg
}

func GetParticipantRecord(letterID uint, userID uint) (*models.LetterParticipant, error) {
	var p models.LetterParticipant
	err := database.DB.Where("letter_id = ? AND user_id = ?", letterID, userID).First(&p).Error
	return &p, err
}

func UpdateLastRead(letterID uint, userID uint) {
	now := time.Now()
	database.DB.Model(&models.LetterParticipant{}).Where("letter_id = ? AND user_id = ?", letterID, userID).Update("last_read_at", now)
}

func SendMessage(letterID uint, senderID uint, body string, attachmentRefs []string) (*models.LetterMessage, error) {
	var msg models.LetterMessage
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		msg = models.LetterMessage{
			LetterID: letterID,
			SenderID: &senderID,
			Body:     body,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		if len(attachmentRefs) > 0 {
			tx.Model(&models.LetterAttachment{}).Where("ref IN ? AND uploader_id = ? AND message_id IS NULL", attachmentRefs, senderID).Update("message_id", msg.ID)
		}

		tx.Model(&models.Letter{}).Where("id = ?", letterID).Update("updated_at", time.Now())

		return nil
	})

	if err == nil {
		database.DB.Preload("Sender").Preload("Attachments").First(&msg, msg.ID)
	}

	return &msg, err
}

func FindMessageByRef(ref string) (*models.LetterMessage, error) {
	var msg models.LetterMessage
	err := database.DB.Preload("Sender").Preload("Attachments").Where("ref = ?", ref).First(&msg).Error
	return &msg, err
}

func EditMessage(msg *models.LetterMessage, body string) error {
	now := time.Now()
	msg.Body = body
	msg.EditedAt = &now
	return database.DB.Save(msg).Error
}

func DeleteMessage(msg *models.LetterMessage) error {
	for _, a := range msg.Attachments {
		database.DB.Delete(&a)
	}
	return database.DB.Delete(msg).Error
}

func RenameLetter(letter *models.Letter, title string) error {
	letter.Title = title
	return database.DB.Save(letter).Error
}

func LeaveLetter(letterID uint, userID uint) error {
	return database.DB.Where("letter_id = ? AND user_id = ?", letterID, userID).Delete(&models.LetterParticipant{}).Error
}

func RemoveParticipant(letterID uint, userID uint) error {
	return database.DB.Where("letter_id = ? AND user_id = ?", letterID, userID).Delete(&models.LetterParticipant{}).Error
}

func GetLetterParticipantCount(letterID uint) int64 {
	var count int64
	database.DB.Model(&models.LetterParticipant{}).Where("letter_id = ?", letterID).Count(&count)
	return count
}

func IsLetterOwner(letterID uint, userID uint) bool {
	var count int64
	database.DB.Model(&models.LetterParticipant{}).Where("letter_id = ? AND user_id = ? AND role = ?", letterID, userID, "owner").Count(&count)
	return count > 0
}

func UpdateAttachmentPath(attachment *models.LetterAttachment) error {
	return database.DB.Save(attachment).Error
}

func FindExistingDM(userID uint, recipientID uint) *models.Letter {
	var participant models.LetterParticipant
	subQuery := database.DB.Model(&models.LetterParticipant{}).Select("letter_id").Where("user_id = ?", recipientID)
	err := database.DB.Where("user_id = ? AND letter_id IN (?)", userID, subQuery).First(&participant).Error
	if err != nil {
		return nil
	}

	count := GetLetterParticipantCount(participant.LetterID)
	if count != 2 {
		return nil
	}

	var letter models.Letter
	if err := database.DB.Where("id = ? AND is_system = ?", participant.LetterID, false).First(&letter).Error; err != nil {
		return nil
	}

	return &letter
}

func UploadAttachment(uploaderID uint, attachment *models.LetterAttachment) error {
	attachment.UploaderID = uploaderID
	return database.DB.Create(attachment).Error
}

func DeleteOrphanedAttachments() {
	var attachments []models.LetterAttachment
	database.DB.Where("message_id IS NULL AND created_at < ?", time.Now().Add(-24*time.Hour)).Find(&attachments)
	for _, a := range attachments {
		database.DB.Delete(&a)
	}
}