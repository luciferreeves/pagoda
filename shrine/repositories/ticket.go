package repositories

import (
	"shrine/database"
	"shrine/models"
	"shrine/utils/meta"
)

func CreateTicketCategory(category *models.TicketCategory) error {
	return database.DB.Create(category).Error
}

func UpdateTicketCategory(category *models.TicketCategory) error {
	return database.DB.Save(category).Error
}

func DeleteTicketCategory(category *models.TicketCategory) error {
	return database.DB.Delete(category).Error
}

func ListTicketCategories() []models.TicketCategory {
	var categories []models.TicketCategory
	database.DB.Order("sort_order asc").Find(&categories)
	return categories
}

func FindTicketCategoryByRef(ref string) (*models.TicketCategory, error) {
	var category models.TicketCategory
	err := database.DB.Where("ref = ?", ref).First(&category).Error
	return &category, err
}

func CreateTicket(ticket *models.Ticket, body string) error {
	if err := database.DB.Create(ticket).Error; err != nil {
		return err
	}
	msg := models.TicketMessage{
		TicketID: ticket.ID,
		SenderID: ticket.UserID,
		Body:     body,
	}
	return database.DB.Create(&msg).Error
}

func FindTicketByRef(ref string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Preload("User").Preload("Category").Preload("Assignee").Where("ref = ?", ref).First(&ticket).Error
	return &ticket, err
}

func ListTicketsForUser(userID uint, p meta.Pagination) ([]models.Ticket, int64) {
	var tickets []models.Ticket
	var total int64

	query := database.DB.Model(&models.Ticket{}).Where("user_id = ?", userID)
	query.Count(&total)
	p.Apply(query.Order("updated_at desc")).Preload("User").Preload("Category").Preload("Assignee").Find(&tickets)

	return tickets, total
}

func ListAllTickets(p meta.Pagination, status string, priority string) ([]models.Ticket, int64) {
	var tickets []models.Ticket
	var total int64

	query := database.DB.Model(&models.Ticket{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	query.Count(&total)
	p.Apply(query.Order("updated_at desc")).Preload("User").Preload("Category").Preload("Assignee").Find(&tickets)

	return tickets, total
}

func GetTicketMessages(ticketID uint) []models.TicketMessage {
	var messages []models.TicketMessage
	database.DB.Where("ticket_id = ?", ticketID).Preload("Sender").Order("created_at asc").Find(&messages)
	return messages
}

func CreateTicketMessage(msg *models.TicketMessage) error {
	return database.DB.Create(msg).Error
}

func UpdateTicket(ticket *models.Ticket) error {
	return database.DB.Save(ticket).Error
}