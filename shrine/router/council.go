package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("council")

	urls.Path(enums.GET, "/users", auth.RequireStaff(controllers.ListUsersController), "users")
	urls.Path(enums.GET, "/users/:username", auth.RequireStaff(controllers.GetUserController), "user")
	urls.Path(enums.POST, "/users/:username/ban", auth.RequireStaff(controllers.BanUserController), "ban")
	urls.Path(enums.POST, "/users/:username/unban", auth.RequireStaff(controllers.UnbanUserController), "unban")
	urls.Path(enums.POST, "/users/:username/disable", auth.RequireStaff(controllers.DisableUserController), "disable")
	urls.Path(enums.POST, "/users/:username/enable", auth.RequireStaff(controllers.EnableUserController), "enable")
	urls.Path(enums.POST, "/users/:username/role", auth.RequireAdmin(controllers.ChangeRoleController), "role")
	urls.Path(enums.PATCH, "/users/:username", auth.RequireAdmin(controllers.EditUserController), "edit")

	urls.Path(enums.GET, "/users/:username/warnings", auth.RequireStaff(controllers.ListWarningsController), "warnings")
	urls.Path(enums.POST, "/users/:username/warn", auth.RequireStaff(controllers.WarnUserController), "warn")
	urls.Path(enums.POST, "/warnings/:ref/deactivate", auth.RequireStaff(controllers.DeactivateWarningController), "deactivate")

	urls.Path(enums.GET, "/tickets", auth.RequireStaff(controllers.ListAllTicketsController), "tickets")
	urls.Path(enums.GET, "/tickets/:ref", auth.RequireStaff(controllers.GetStaffTicketController), "ticket")
	urls.Path(enums.PATCH, "/tickets/:ref", auth.RequireStaff(controllers.UpdateTicketController), "ticketupdate")
	urls.Path(enums.POST, "/tickets/:ref/messages", auth.RequireStaff(controllers.StaffReplyTicketController), "ticketreply")
	urls.Path(enums.GET, "/categories", auth.RequireStaff(controllers.ListTicketCategoriesController), "categories")
	urls.Path(enums.POST, "/categories", auth.RequireAdmin(controllers.CreateTicketCategoryController), "categorycreate")
	urls.Path(enums.PATCH, "/categories/:ref", auth.RequireAdmin(controllers.UpdateTicketCategoryController), "categoryupdate")
	urls.Path(enums.DELETE, "/categories/:ref", auth.RequireAdmin(controllers.DeleteTicketCategoryController), "categorydelete")

	urls.Path(enums.GET, "/audit", auth.RequireStaff(controllers.ListAuditLogsController), "audit")
	urls.Path(enums.GET, "/audit/:ref", auth.RequireStaff(controllers.GetAuditLogController), "auditdetail")

	urls.Path(enums.GET, "/bannedips", auth.RequireAdmin(controllers.ListIPBansController), "bannedips")
	urls.Path(enums.DELETE, "/bannedips/:id", auth.RequireAdmin(controllers.DeleteIPBanController), "bannedipdelete")

	urls.Path(enums.POST, "/upload", auth.RequireAdmin(controllers.UploadImageController), "upload")

	urls.Path(enums.GET, "/districts/requests", auth.RequireStaff(controllers.ListSiteRequestsController), "districtreqs")
	urls.Path(enums.POST, "/districts/sites/:ref/review", auth.RequireStaff(controllers.ReviewSiteController), "districtreview")
	urls.Path(enums.GET, "/districts/sites", auth.RequireStaff(controllers.ListAdminSitesController), "districtsites")
	urls.Path(enums.PATCH, "/districts/sites/:ref", auth.RequireStaff(controllers.EditSiteController), "districtedit")
	urls.Path(enums.GET, "/districts/pending", auth.RequireStaff(controllers.CountPendingSitesController), "districtpending")
}