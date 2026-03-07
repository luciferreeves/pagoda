package services

import (
	"shrine/repositories"
	"shrine/types/user"
)

func GetStats() user.StatsResponse {
	return user.StatsResponse{
		Citizens:       repositories.CountCitizens(),
		Online:         repositories.CountOnline(),
		NewestCitizens: buildCitizenSummaries(repositories.NewestCitizens(5)),
		OnlineCitizens: buildCitizenSummaries(repositories.OnlineCitizens(10)),
	}
}