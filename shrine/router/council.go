package router

import (
	"shrine/controllers"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("council")

	urls.Path(types.GET, "/users", auth.RequireStaff(controllers.ListUsersController), "users")
	urls.Path(types.GET, "/users/:username", auth.RequireStaff(controllers.GetUserController), "user")
	urls.Path(types.POST, "/users/:username/ban", auth.RequireStaff(controllers.BanUserController), "ban")
	urls.Path(types.POST, "/users/:username/unban", auth.RequireStaff(controllers.UnbanUserController), "unban")
	urls.Path(types.POST, "/users/:username/disable", auth.RequireStaff(controllers.DisableUserController), "disable")
	urls.Path(types.POST, "/users/:username/enable", auth.RequireStaff(controllers.EnableUserController), "enable")
	urls.Path(types.POST, "/users/:username/role", auth.RequireAdmin(controllers.ChangeRoleController), "role")
}