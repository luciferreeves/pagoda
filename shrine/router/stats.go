package router

import (
	"shrine/controllers"
	"shrine/types"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("stats")

	urls.Path(types.GET, "/", controllers.StatsController, "index")
}