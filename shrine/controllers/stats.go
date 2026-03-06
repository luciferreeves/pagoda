package controllers

import (
	"shrine/repositories"
	"shrine/types"

	"github.com/gofiber/fiber/v2"
)

func StatsController(context *fiber.Ctx) error {
	newest := repositories.NewestCitizens(5)
	online := repositories.OnlineCitizens(10)

	newestSummaries := make([]types.CitizenSummary, len(newest))
	for i, u := range newest {
		newestSummaries[i] = u.ToSummary()
	}

	onlineSummaries := make([]types.CitizenSummary, len(online))
	for i, u := range online {
		onlineSummaries[i] = u.ToSummary()
	}

	return Success(context, types.StatsResponse{
		Citizens:       repositories.CountCitizens(),
		Online:         repositories.CountOnline(),
		NewestCitizens: newestSummaries,
		OnlineCitizens: onlineSummaries,
	})
}