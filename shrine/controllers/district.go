package controllers

import (
	"shrine/services"
	"shrine/types/district"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ListDistrictsController(context *fiber.Ctx) error {
	return shortcuts.Success(context, services.ListDistricts())
}

func ListDistrictSitesController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	request := meta.Request(context)
	slug, _ := request.Query("district")
	tag, _ := request.Query("tag")
	search, _ := request.Query("search")

	items, total := services.ListDistrictSites(pagination, slug, tag, search)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func SubmitSiteController(context *fiber.Ctx) error {
	citizen := auth.GetUser(context)

	body, err := meta.Body[district.SubmitSiteRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.SubmitSite(citizen.ID, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Created(context, result)
}

func ListSiteRequestsController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	status, _ := meta.Request(context).Query("status")

	items, total := services.ListSiteRequests(pagination, status)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func ReviewSiteController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[district.ReviewSiteRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.ReviewSite(admin.ID, ref, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func ListAdminSitesController(context *fiber.Ctx) error {
	pagination := meta.Paginate(context)
	request := meta.Request(context)
	slug, _ := request.Query("district")
	search, _ := request.Query("search")

	items, total := services.ListAdminSites(pagination, slug, search)
	return shortcuts.Success(context, pagination.Response(items, total))
}

func EditSiteController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)
	ref := meta.Request(context).MustHave().Param("ref")

	body, err := meta.Body[district.EditSiteRequest](context)
	if err != nil {
		return shortcuts.BadRequest(context, err)
	}

	result, serviceErr := services.EditSite(admin.ID, ref, body)
	if serviceErr != nil {
		return shortcuts.HandleError(context, serviceErr)
	}

	return shortcuts.Success(context, result)
}

func CountPendingSitesController(context *fiber.Ctx) error {
	return shortcuts.Success(context, fiber.Map{"count": services.CountPendingSites()})
}