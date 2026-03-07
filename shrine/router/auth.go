package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(enums.POST, "/register", controllers.RegisterController, "register")
	urls.Path(enums.POST, "/login", controllers.LoginController, "login")
	urls.Path(enums.POST, "/verify", controllers.VerifyController, "verify")
	urls.Path(enums.POST, "/reactivate", controllers.ResendActivationController, "reactivate")
	urls.Path(enums.POST, "/logout", auth.RequireAuthentication(controllers.LogoutController), "logout")
	urls.Path(enums.GET, "/me", auth.RequireAuthentication(controllers.MeController), "me")
	urls.Path(enums.POST, "/heartbeat", auth.RequireAuthentication(controllers.HeartbeatController), "heartbeat")
}