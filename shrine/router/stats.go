package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("stats")

	urls.Path(enums.GET, "/", controllers.StatsController, "index")
}