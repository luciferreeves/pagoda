package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("config")

	urls.Path(enums.GET, "/", controllers.ConfigController, "index")
}