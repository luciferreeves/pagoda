package services

import (
	"fmt"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/audit"
	"shrine/types/hypertext"
	"shrine/types/warning"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/sanitize"
	"strings"
)

func WarnUser(admin *models.User, target *models.User, request warning.WarnRequest) (*warning.WarningResponse, *hypertext.ServiceError) {
	if hierarchyErr := auth.ValidateHierarchy(admin, target, "warn"); hierarchyErr != nil {
		return nil, fail(enums.Forbidden, hierarchyErr.Error())
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		return nil, fail(enums.BadRequest, messages.WarningTitleRequired)
	}

	sanitizedMessage := sanitize.HTML(request.Message)
	if sanitizedMessage == "" {
		return nil, fail(enums.BadRequest, messages.WarningMessageRequired)
	}

	record, err := repositories.CreateWarning(admin.ID, target.ID, title, sanitizedMessage)
	if err != nil {
		return nil, fail(enums.Internal, messages.FailedCreateWarning)
	}

	repositories.LogAction(admin.ID, "user.warn", "user", target.Username, fmt.Sprintf("Warned user @%s", target.Username), audit.WarningDetails{
		WarningRef: record.SystemRef,
		Title:      title,
		Message:    sanitizedMessage,
	})

	response := record.ToResponse()
	return &response, nil
}

func DeactivateWarning(admin *models.User, ref string) (*warning.WarningResponse, *hypertext.ServiceError) {
	record, err := repositories.FindWarningByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.WarningNotFound)
	}

	if !record.Active {
		return nil, fail(enums.BadRequest, messages.WarningAlreadyInactive)
	}

	if err := repositories.DeactivateWarning(record); err != nil {
		return nil, fail(enums.Internal, messages.FailedDeactivateWarn)
	}

	repositories.LogAction(admin.ID, "user.unwarn", "user", "", fmt.Sprintf("Deactivated warning %s", record.SystemRef), audit.DeactivateWarningDetails{
		WarningRef: record.SystemRef,
	})

	response := record.ToResponse()
	return &response, nil
}

func ListWarnings(username string, pagination meta.Pagination) ([]warning.WarningResponse, int64, *hypertext.ServiceError) {
	citizen, serviceErr := ResolveUser(username)
	if serviceErr != nil {
		return nil, 0, serviceErr
	}

	warnings, total := repositories.ListWarningsForUser(citizen.ID, pagination)

	items := make([]warning.WarningResponse, len(warnings))
	for index, record := range warnings {
		items[index] = record.ToResponse()
	}

	return items, total, nil
}