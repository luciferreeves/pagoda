package repositories

import (
	"shrine/database"
	"shrine/models"
	"shrine/utils/meta"
)

func CreateIPBan(ip string, reason string) error {
	return database.DB.Create(&models.IPBan{
		IP:     ip,
		Reason: reason,
	}).Error
}

func DeleteIPBan(ip string) {
	database.DB.Where("ip = ?", ip).Delete(&models.IPBan{})
}

func IsIPBanned(ip string) bool {
	var count int64
	database.DB.Model(&models.IPBan{}).Where("ip = ?", ip).Count(&count)
	return count > 0
}

func ListIPBans(pagination meta.Pagination, sorting meta.Sorting) ([]models.IPBan, int64) {
	var ipBans []models.IPBan
	var total int64

	query := database.DB.Model(&models.IPBan{})
	query.Count(&total)
	pagination.Apply(sorting.Apply(query)).Find(&ipBans)

	return ipBans, total
}

func DeleteIPBanByID(id uint) {
	database.DB.Delete(&models.IPBan{}, id)
}