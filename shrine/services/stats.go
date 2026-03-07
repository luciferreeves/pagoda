package services

import (
	"shrine/repositories"
	"shrine/types/user"
)

func GetStats() user.StatsResponse {
	newest := repositories.NewestCitizens(5)
	online := repositories.OnlineCitizens(10)

	newestSummaries := make([]user.CitizenSummaryResponse, len(newest))
	for index, citizen := range newest {
		newestSummaries[index] = citizen.ToSummary()
	}

	onlineSummaries := make([]user.CitizenSummaryResponse, len(online))
	for index, citizen := range online {
		onlineSummaries[index] = citizen.ToSummary()
	}

	return user.StatsResponse{
		Citizens:       repositories.CountCitizens(),
		Online:         repositories.CountOnline(),
		NewestCitizens: newestSummaries,
		OnlineCitizens: onlineSummaries,
	}
}