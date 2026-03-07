package services

import (
	"shrine/enums"
	"shrine/messages"
	"shrine/repositories"
	"shrine/types/audit"
	"shrine/types/hypertext"
	"shrine/utils/meta"
)

func ListAuditLogs(pagination meta.Pagination, action string, targetType string) ([]audit.AuditLogResponse, int64) {
	logs, total := repositories.ListAuditLogs(pagination, action, targetType)

	items := make([]audit.AuditLogResponse, len(logs))
	for index, record := range logs {
		items[index] = record.ToResponse()
	}

	return items, total
}

func GetAuditLog(ref string) (*audit.DetailResponse, *hypertext.ServiceError) {
	record, err := repositories.FindAuditLogByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.AuditLogNotFound)
	}

	response := record.ToDetailResponse()
	return &response, nil
}