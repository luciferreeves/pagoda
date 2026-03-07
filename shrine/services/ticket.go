package services

import (
	"fmt"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/common"
	"shrine/types/hypertext"
	"shrine/types/ticket"
	"shrine/utils/meta"
	"strings"
)

func ListUserTickets(userID uint, pagination meta.Pagination) ([]ticket.TicketResponse, int64) {
	tickets, total := repositories.ListTicketsForUser(userID, pagination)

	items := make([]ticket.TicketResponse, len(tickets))
	for index, record := range tickets {
		items[index] = record.ToResponse()
	}

	return items, total
}

func GetUserTicket(ref string, userID uint) (*ticket.DetailResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveTicket(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if record.UserID != userID {
		return nil, fail(enums.NotFound, messages.TicketNotFound)
	}

	ticketMessages := repositories.GetTicketMessages(record.ID)

	response := ticket.DetailResponse{
		TicketResponse: record.ToResponse(),
		Messages:       buildTicketMessageResponses(ticketMessages),
	}

	return &response, nil
}

func CreateTicket(userID uint, request ticket.CreateRequest) (*common.MessageResponse, *hypertext.ServiceError) {
	category, err := repositories.FindTicketCategoryByRef(request.CategoryRef)
	if err != nil {
		return nil, fail(enums.BadRequest, messages.InvalidCategory)
	}

	subject := strings.TrimSpace(request.Subject)
	if subject == "" || len(subject) > 200 {
		return nil, fail(enums.BadRequest, messages.SubjectRequired)
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	priority := enums.TicketPriority(request.Priority)
	if priority == "" {
		priority = enums.PriorityLow
	}

	switch priority {
	case enums.PriorityLow, enums.PriorityMedium, enums.PriorityHigh, enums.PriorityUrgent:
	default:
		return nil, fail(enums.BadRequest, messages.InvalidPriority)
	}

	record := models.Ticket{
		UserID:     userID,
		CategoryID: category.ID,
		Subject:    subject,
		Priority:   string(priority),
	}

	if err := repositories.CreateTicket(&record, sanitizedBody); err != nil {
		return nil, fail(enums.Internal, messages.FailedCreateTicket)
	}

	return &common.MessageResponse{Message: record.Ref}, nil
}

func ReplyTicket(ref string, userID uint, request ticket.SendMessageRequest) (*ticket.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveTicket(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if record.UserID != userID {
		return nil, fail(enums.NotFound, messages.TicketNotFound)
	}

	if record.Status == string(enums.StatusClosed) {
		return nil, fail(enums.BadRequest, messages.TicketClosed)
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	message := models.TicketMessage{
		TicketID: record.ID,
		SenderID: userID,
		Body:     sanitizedBody,
	}

	if err := repositories.CreateTicketMessage(&message); err != nil {
		return nil, fail(enums.Internal, messages.FailedSendReply)
	}

	response := message.ToResponse()
	return &response, nil
}

func ListAllTickets(pagination meta.Pagination, status string, priority string) ([]ticket.TicketResponse, int64) {
	tickets, total := repositories.ListAllTickets(pagination, status, priority)

	items := make([]ticket.TicketResponse, len(tickets))
	for index, record := range tickets {
		items[index] = record.ToResponse()
	}

	return items, total
}

func GetStaffTicket(ref string) (*ticket.DetailResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveTicket(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	ticketMessages := repositories.GetTicketMessages(record.ID)

	response := ticket.DetailResponse{
		TicketResponse: record.ToResponse(),
		Messages:       buildTicketMessageResponses(ticketMessages),
	}

	return &response, nil
}

func UpdateTicket(adminID uint, ref string, request ticket.UpdateRequest) (*ticket.TicketResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveTicket(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if request.Status != nil {
		status := enums.TicketStatus(*request.Status)
		switch status {
		case enums.StatusOpen, enums.StatusInProgress, enums.StatusResolved, enums.StatusClosed:
			record.Status = *request.Status
		default:
			return nil, fail(enums.BadRequest, messages.InvalidStatus)
		}
	}

	if request.Priority != nil {
		priority := enums.TicketPriority(*request.Priority)
		switch priority {
		case enums.PriorityLow, enums.PriorityMedium, enums.PriorityHigh, enums.PriorityUrgent:
			record.Priority = *request.Priority
		default:
			return nil, fail(enums.BadRequest, messages.InvalidPriority)
		}
	}

	if request.CategoryRef != nil {
		category, err := repositories.FindTicketCategoryByRef(*request.CategoryRef)
		if err != nil {
			return nil, fail(enums.BadRequest, messages.InvalidCategory)
		}
		record.CategoryID = category.ID
	}

	if request.Assignee != nil {
		if *request.Assignee == "" {
			record.AssigneeID = nil
		} else {
			assignee, serviceErr := ResolveUser(*request.Assignee)
			if serviceErr != nil {
				return nil, fail(enums.BadRequest, messages.AssigneeNotFound)
			}
			record.AssigneeID = &assignee.ID
		}
	}

	if err := repositories.UpdateTicket(record); err != nil {
		return nil, fail(enums.Internal, messages.FailedUpdateTicket)
	}

	repositories.LogAction(adminID, "ticket.update", "ticket", record.Ref, fmt.Sprintf("Updated ticket %s", record.Ref), request)

	record, _ = repositories.FindTicketByRef(record.Ref)
	response := record.ToResponse()
	return &response, nil
}

func StaffReplyTicket(ref string, staffID uint, request ticket.SendMessageRequest) (*ticket.MessageResponse, *hypertext.ServiceError) {
	record, serviceErr := resolveTicket(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	sanitizedBody, serviceErr := sanitizeRequiredBody(request.Body)
	if serviceErr != nil {
		return nil, serviceErr
	}

	message := models.TicketMessage{
		TicketID: record.ID,
		SenderID: staffID,
		Body:     sanitizedBody,
		IsStaff:  true,
	}

	if err := repositories.CreateTicketMessage(&message); err != nil {
		return nil, fail(enums.Internal, messages.FailedSendReply)
	}

	response := message.ToResponse()
	return &response, nil
}

func ListTicketCategories() []ticket.CategoryResponse {
	categories := repositories.ListTicketCategories()

	items := make([]ticket.CategoryResponse, len(categories))
	for index, category := range categories {
		items[index] = category.ToResponse()
	}

	return items
}

func CreateTicketCategory(request ticket.CreateCategoryRequest) (*ticket.CategoryResponse, *hypertext.ServiceError) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return nil, fail(enums.BadRequest, messages.CategoryNameRequired)
	}

	category := models.TicketCategory{
		Name:        name,
		Description: strings.TrimSpace(request.Description),
	}

	if err := repositories.CreateTicketCategory(&category); err != nil {
		return nil, fail(enums.Internal, messages.FailedCreateCategory)
	}

	response := category.ToResponse()
	return &response, nil
}

func UpdateTicketCategory(ref string, request ticket.UpdateCategoryRequest) (*ticket.CategoryResponse, *hypertext.ServiceError) {
	category, err := repositories.FindTicketCategoryByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.CategoryNotFound)
	}

	if request.Name != nil {
		category.Name = strings.TrimSpace(*request.Name)
	}
	if request.Description != nil {
		category.Description = strings.TrimSpace(*request.Description)
	}
	if request.SortOrder != nil {
		category.SortOrder = *request.SortOrder
	}

	if err := repositories.UpdateTicketCategory(category); err != nil {
		return nil, fail(enums.Internal, messages.FailedUpdateCategory)
	}

	response := category.ToResponse()
	return &response, nil
}

func DeleteTicketCategory(ref string) *hypertext.ServiceError {
	category, err := repositories.FindTicketCategoryByRef(ref)
	if err != nil {
		return fail(enums.NotFound, messages.CategoryNotFound)
	}

	if err := repositories.DeleteTicketCategory(category); err != nil {
		return fail(enums.Internal, messages.FailedDeleteCategory)
	}

	return nil
}