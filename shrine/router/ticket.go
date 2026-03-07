package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("tickets")

	urls.Path(enums.GET, "", auth.RequireAuthentication(controllers.ListUserTicketsController), "list")
	urls.Path(enums.POST, "", auth.RequireAuthentication(controllers.CreateTicketController), "create")
	urls.Path(enums.GET, "/:ref", auth.RequireAuthentication(controllers.GetUserTicketController), "detail")
	urls.Path(enums.POST, "/:ref/messages", auth.RequireAuthentication(controllers.ReplyTicketController), "reply")
}