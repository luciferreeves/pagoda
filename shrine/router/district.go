package router

import (
	"shrine/controllers"
	"shrine/enums"
	"shrine/utils/auth"
	"shrine/utils/urls"
)

func init() {
	urls.SetNamespace("districts")

	urls.Path(enums.GET, "", controllers.ListDistrictsController, "list")
	urls.Path(enums.GET, "/sites", controllers.ListDistrictSitesController, "sites")
	urls.Path(enums.POST, "/sites", auth.RequireAuthentication(controllers.SubmitSiteController), "submit")
}