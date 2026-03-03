package router

import (
	"shrine/controllers"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(types.POST, "/register", controllers.RegisterController, "register")
	urls.Path(types.POST, "/login", controllers.LoginController, "login")
	urls.Path(types.POST, "/logout", auth.RequireAuthentication(controllers.LogoutController), "logout")
	urls.Path(types.GET, "/me", auth.RequireAuthentication(controllers.MeController), "me")
}