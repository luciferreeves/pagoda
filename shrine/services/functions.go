package services

import (
	"fmt"
	"regexp"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/hypertext"
	"shrine/types/letter"
	"shrine/types/ticket"
	"shrine/types/user"
	"shrine/utils/sanitize"
	"strings"
)

func ResolveUser(username string) (*models.User, *hypertext.ServiceError) {
	citizen, err := repositories.FindUserByUsername(username)
	if err != nil {
		return nil, fail(enums.NotFound, messages.UserNotFound)
	}
	return citizen, nil
}

func mapRegistrationError(err error) *hypertext.ServiceError {
	return mapUserError(err)
}

func mapUserError(err error) *hypertext.ServiceError {
	if strings.Contains(err.Error(), "users.username") {
		return fail(enums.Conflict, messages.UsernameAlreadyExists)
	}
	if strings.Contains(err.Error(), "users.email") {
		return fail(enums.Conflict, messages.EmailAlreadyExists)
	}
	return fail(enums.BadRequest, err.Error())
}

func assembleLetterResponse(record *models.Letter, viewerID uint) letter.LetterResponse {
	participants := repositories.GetLetterParticipants(record.ID)
	lastMessage := repositories.GetLastMessage(record.ID)
	viewer, _ := repositories.GetParticipantRecord(record.ID, viewerID)

	var unread bool
	if viewer != nil && lastMessage != nil {
		unread = viewer.LastReadAt == nil || lastMessage.CreatedAt.After(*viewer.LastReadAt)
	}

	var lastMessageResponse *letter.MessageResponse
	if lastMessage != nil {
		response := lastMessage.ToResponse()
		lastMessageResponse = &response
	}

	return letter.LetterResponse{
		Ref:          record.Ref,
		Title:        computeTitle(record, participants, viewerID),
		IsSystem:     record.IsSystem,
		SystemRef:    record.SystemRef,
		Participants: buildParticipantResponses(participants),
		LastMessage:  lastMessageResponse,
		Unread:       unread,
		UpdatedAt:    record.UpdatedAt,
	}
}

func computeTitle(record *models.Letter, participants []models.LetterParticipant, viewerID uint) string {
	if record.Title != "" {
		return record.Title
	}

	if record.IsSystem {
		return messages.SystemMessageTitle
	}

	var others []string
	for _, participant := range participants {
		if participant.UserID != viewerID {
			others = append(others, participant.User.DisplayName)
		}
	}

	switch len(others) {
	case 0:
		return messages.EmptyConversationTitle
	case 1:
		return others[0]
	case 2:
		return fmt.Sprintf(messages.LetterTitleTwo, others[0], others[1])
	default:
		return fmt.Sprintf(messages.LetterTitleMany, others[0], others[1], len(others)-2)
	}
}

func buildParticipantResponses(participants []models.LetterParticipant) []letter.ParticipantResponse {
	responses := make([]letter.ParticipantResponse, len(participants))
	for index, participant := range participants {
		responses[index] = participant.ToResponse()
	}
	return responses
}

func buildMessageResponses(letterMessages []models.LetterMessage) []letter.MessageResponse {
	responses := make([]letter.MessageResponse, len(letterMessages))
	for index, message := range letterMessages {
		responses[index] = message.ToResponse()
	}
	return responses
}

func resolveLetter(ref string, userID uint) (*models.Letter, *hypertext.ServiceError) {
	record, err := repositories.FindLetterByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.LetterNotFound)
	}

	if !repositories.IsLetterParticipant(record.ID, userID) {
		return nil, fail(enums.NotFound, messages.LetterNotFound)
	}

	return record, nil
}

func resolveTicket(ref string) (*models.Ticket, *hypertext.ServiceError) {
	record, err := repositories.FindTicketByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.TicketNotFound)
	}
	return record, nil
}

func sanitizeRequiredBody(body string) (string, *hypertext.ServiceError) {
	sanitized := sanitize.HTML(body)
	if sanitized == "" {
		return "", fail(enums.BadRequest, messages.MessageBodyRequired)
	}
	return sanitized, nil
}

func buildTicketMessageResponses(ticketMessages []models.TicketMessage) []ticket.MessageResponse {
	responses := make([]ticket.MessageResponse, len(ticketMessages))
	for index, message := range ticketMessages {
		responses[index] = message.ToResponse()
	}
	return responses
}

func buildCitizenSummaries(citizens []models.User) []user.CitizenSummaryResponse {
	summaries := make([]user.CitizenSummaryResponse, len(citizens))
	for index, citizen := range citizens {
		summaries[index] = citizen.ToSummary()
	}
	return summaries
}

var imgSrcPattern = regexp.MustCompile(`<img\s+[^>]*src="([^"]+)"`)

func extractImageURL(html string) string {
	match := imgSrcPattern.FindStringSubmatch(html)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}