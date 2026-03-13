package services

import (
	"shrine/models"
	"shrine/repositories"
	"shrine/types/user"
)

func GetStats(citizen *models.User) user.StatsResponse {
	response := user.StatsResponse{
		Citizens:       repositories.CountCitizens(),
		Online:         repositories.CountOnline(),
		NewestCitizens: buildCitizenSummaries(repositories.NewestCitizens(5)),
		OnlineCitizens: buildCitizenSummaries(repositories.OnlineCitizens(10)),
	}

	if citizen != nil {
		response.UnreadLetters = repositories.CountUnreadLetters(citizen.ID)
		if citizen.IsStaff() {
			response.PendingDistricts = repositories.CountPendingDistrictSites()
		}
	}

	return response
}