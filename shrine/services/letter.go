package services

import (
	"fmt"
	"io"
	"shrine/config"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/common"
	"shrine/types/hypertext"
	"shrine/types/letter"
	"shrine/utils/files"
	"shrine/utils/meta"
	"shrine/utils/storage"
	"strings"
)

func ListLetters(userID uint, pagination meta.Pagination) ([]letter.LetterResponse, int64) {
	letters, total := repositories.ListLettersForUser(userID, pagination)

	items := make([]letter.LetterResponse, len(letters))
	for index, record := range letters {
		items[index] = assembleLetterResponse(&record, userID)
	}

	return items, total
}

func GetLetter(ref string, userID uint, pagination meta.Pagination) (*letter.DetailResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveLetter(ref, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	participants := repositories.GetLetterParticipants(record.ID)
	letterMessages, _ := repositories.GetLetterMessages(record.ID, pagination)

	repositories.UpdateLastRead(record.ID, userID)

	response := letter.DetailResponse{
		Ref:          record.Ref,
		Title:        computeTitle(record, participants, userID),
		IsSystem:     record.IsSystem,
		SystemRef:    record.SystemRef,
		Participants: buildParticipantResponses(participants),
		Messages:     buildMessageResponses(letterMessages),
		CreatedAt:    record.CreatedAt,
	}

	return &response, nil
}

func CreateLetter(userID uint, request letter.CreateRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	if len(request.Recipients) == 0 {
		return nil, fail(enums.BadRequest, messages.RecipientsRequired)
	}

	if len(request.Recipients) > 20 {
		return nil, fail(enums.BadRequest, messages.RecipientsMax)
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	recipientIDs := make([]uint, 0, len(request.Recipients))
	for _, username := range request.Recipients {
		recipient, err := repositories.FindUserByUsername(username)
		if err != nil {
			return nil, fail(enums.BadRequest, fmt.Sprintf(messages.RecipientNotFound, username))
		}
		if recipient.ID == userID {
			continue
		}
		if recipient.LetterPrivacy == enums.LetterPrivacyFriends {
			return nil, fail(enums.BadRequest, fmt.Sprintf(messages.RecipientNotAcceptingLetters, username))
		}
		recipientIDs = append(recipientIDs, recipient.ID)
	}

	if len(recipientIDs) == 0 {
		return nil, fail(enums.BadRequest, messages.ValidRecipientRequired)
	}

	if len(recipientIDs) == 1 && request.Title == "" {
		existing := repositories.FindExistingDM(userID, recipientIDs[0])
		if existing != nil {
			_, err := repositories.SendMessage(existing.ID, userID, sanitizedBody, request.AttachmentRefs)
			if err != nil {
				return nil, fail(enums.Internal, messages.FailedSendMessage)
			}
			return &common.MessageResponse{Message: existing.Ref}, nil
		}
	}

	record, err := repositories.CreateLetter(userID, strings.TrimSpace(request.Title), recipientIDs, sanitizedBody, request.AttachmentRefs)
	if err != nil {
		return nil, fail(enums.Internal, messages.FailedCreateLetter)
	}

	return &common.MessageResponse{Message: record.Ref}, nil
}

func SendLetterMessage(ref string, userID uint, request letter.SendMessageRequest) (*letter.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveLetter(ref, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if record.IsSystem {
		return nil, fail(enums.BadRequest, messages.SystemLetterNotReplyable)
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	message, err := repositories.SendMessage(record.ID, userID, sanitizedBody, request.AttachmentRefs)
	if err != nil {
		return nil, fail(enums.Internal, messages.FailedSendMessage)
	}

	response := message.ToResponse()
	return &response, nil
}

func EditLetterMessage(letterRef string, messageRef string, userID uint, request letter.EditMessageRequest) (*letter.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveLetter(letterRef, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	message, err := repositories.FindMessageByRef(messageRef)
	if err != nil || message.LetterID != record.ID {
		return nil, fail(enums.NotFound, messages.MessageNotFound)
	}

	if message.SenderID == nil || *message.SenderID != userID {
		return nil, fail(enums.Forbidden, messages.CannotEditOthersMessage)
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if err := repositories.EditMessage(message, sanitizedBody); err != nil {
		return nil, fail(enums.Internal, messages.FailedEditMessage)
	}

	response := message.ToResponse()
	return &response, nil
}

func DeleteLetterMessage(letterRef string, messageRef string, userID uint) *hypertext.ServiceError {
	record, serviceErr := resolveLetter(letterRef, userID)
	if serviceErr != nil {
		return serviceErr
	}

	message, err := repositories.FindMessageByRef(messageRef)
	if err != nil || message.LetterID != record.ID {
		return fail(enums.NotFound, messages.MessageNotFound)
	}

	if message.SenderID == nil || *message.SenderID != userID {
		return fail(enums.Forbidden, messages.CannotDeleteOthersMessage)
	}

	for _, attachment := range message.Attachments {
		storage.Delete(attachment.FilePath)
	}

	if err := repositories.DeleteMessage(message); err != nil {
		return fail(enums.Internal, messages.FailedDeleteMessage)
	}

	return nil
}

func RenameLetter(ref string, userID uint, request letter.RenameRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveLetter(ref, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if record.IsSystem {
		return nil, fail(enums.BadRequest, messages.SystemLetterNotRenameable)
	}

	title := strings.TrimSpace(request.Title)
	if len(title) > 200 {
		return nil, fail(enums.BadRequest, messages.TitleTooLong)
	}

	if err := repositories.RenameLetter(record, title); err != nil {
		return nil, fail(enums.Internal, messages.FailedRenameLetter)
	}

	return &common.MessageResponse{Message: messages.LetterRenamed}, nil
}

func LeaveLetter(ref string, userID uint) (*common.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveLetter(ref, userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if record.IsSystem {
		return nil, fail(enums.BadRequest, messages.SystemLetterNotLeaveable)
	}

	if err := repositories.LeaveLetter(record.ID, userID); err != nil {
		return nil, fail(enums.Internal, messages.FailedLeaveLetter)
	}

	return &common.MessageResponse{Message: messages.LeftConversation}, nil
}

func RemoveLetterParticipant(ref string, ownerID uint, request letter.RemoveParticipantRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	record, err := repositories.FindLetterByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.LetterNotFound)
	}

	if !repositories.IsLetterOwner(record.ID, ownerID) {
		return nil, fail(enums.Forbidden, messages.OnlyOwnerCanRemove)
	}

	target, serviceErr := ResolveUser(request.Username)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if target.ID == ownerID {
		return nil, fail(enums.BadRequest, messages.CannotRemoveSelf)
	}

	if !repositories.IsLetterParticipant(record.ID, target.ID) {
		return nil, fail(enums.BadRequest, messages.UserNotInConversation)
	}

	if err := repositories.RemoveParticipant(record.ID, target.ID); err != nil {
		return nil, fail(enums.Internal, messages.FailedRemoveParticipant)
	}

	return &common.MessageResponse{Message: fmt.Sprintf(messages.ParticipantRemoved, target.DisplayName)}, nil
}

func UploadLetterAttachment(userID uint, fileName string, fileSize int64, contentType string, reader io.Reader) (*letter.AttachmentResponse, *hypertext.ServiceError) {
	if fileSize > config.Storage.MaxFileSize {
		return nil, fail(enums.BadRequest, fmt.Sprintf(messages.FileExceedsMaxSize, config.Storage.MaxFileSize/1024/1024))
	}

	attachment := models.LetterAttachment{
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
		Category:    files.DetectCategory(contentType),
	}

	if err := repositories.UploadAttachment(userID, &attachment); err != nil {
		return nil, fail(enums.Internal, messages.FailedSaveAttachment)
	}

	path := fmt.Sprintf("letters/attachments/%s/%s", attachment.Ref, fileName)
	if err := storage.Upload(path, reader, fileSize, contentType); err != nil {
		return nil, fail(enums.Internal, messages.FailedUploadFile)
	}

	attachment.FilePath = path
	if err := repositories.UpdateAttachmentPath(&attachment); err != nil {
		return nil, fail(enums.Internal, messages.FailedSaveAttachment)
	}

	response := attachment.ToResponse()
	return &response, nil
}