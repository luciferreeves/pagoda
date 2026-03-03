package router

import (
	"shrine/controllers"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("sample")

	urls.Path(types.GET, "/hello", auth.RequireAuthentication(controllers.HelloController), "hello")
}
