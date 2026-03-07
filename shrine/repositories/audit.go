package repositories

import (
	"encoding/json"
	"shrine/database"
	"shrine/models"
	"shrine/utils/crypto"
	"shrine/utils/meta"
)

func LogAction(actorID uint, action string, targetType string, targetRef string, summary string, details any) string {
	detailsJSON, _ := json.Marshal(details)
	ref := crypto.SystemRef()

	log := models.AuditLog{
		SystemRef:  ref,
		ActorID:    actorID,
		Action:     action,
		TargetType: targetType,
		TargetRef:  targetRef,
		Summary:    summary,
		Details:    string(detailsJSON),
	}

	database.DB.Create(&log)
	return ref
}

func ListAuditLogs(p meta.Pagination, action string, targetType string) ([]models.AuditLog, int64) {
	var logs []models.AuditLog
	var total int64

	query := database.DB.Model(&models.AuditLog{})

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	query.Count(&total)
	p.Apply(query.Order("created_at desc")).Preload("Actor").Find(&logs)

	return logs, total
}

func FindAuditLogByRef(ref string) (*models.AuditLog, error) {
	var log models.AuditLog
	err := database.DB.Preload("Actor").Where("system_ref = ?", ref).First(&log).Error
	return &log, err
}