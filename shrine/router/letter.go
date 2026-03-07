package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("letters")

	urls.Path(enums.GET, "", auth.RequireAuthentication(controllers.ListLettersController), "list")
	urls.Path(enums.POST, "", auth.RequireAuthentication(controllers.CreateLetterController), "create")
	urls.Path(enums.POST, "/attachments", auth.RequireAuthentication(controllers.UploadAttachmentController), "upload")
	urls.Path(enums.GET, "/:ref", auth.RequireAuthentication(controllers.GetLetterController), "detail")
	urls.Path(enums.PATCH, "/:ref", auth.RequireAuthentication(controllers.RenameLetterController), "rename")
	urls.Path(enums.POST, "/:ref/messages", auth.RequireAuthentication(controllers.SendMessageController), "send")
	urls.Path(enums.PATCH, "/:ref/messages/:messageRef", auth.RequireAuthentication(controllers.EditMessageController), "edit")
	urls.Path(enums.DELETE, "/:ref/messages/:messageRef", auth.RequireAuthentication(controllers.DeleteMessageController), "delete")
	urls.Path(enums.POST, "/:ref/leave", auth.RequireAuthentication(controllers.LeaveLetterController), "leave")
	urls.Path(enums.POST, "/:ref/remove", auth.RequireAuthentication(controllers.RemoveParticipantController), "remove")
}