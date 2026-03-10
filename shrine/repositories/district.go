package repositories

import (
	"shrine/database"
	"shrine/models"
	"shrine/utils/meta"
	"strings"
)

func CreateDistrictSite(site *models.DistrictSite) error {
	return database.DB.Create(site).Error
}

func UpdateDistrictSite(site *models.DistrictSite) error {
	return database.DB.Save(site).Error
}

func FindDistrictSiteByRef(ref string) (*models.DistrictSite, error) {
	var site models.DistrictSite
	err := database.DB.Preload("Submitter").Preload("ReviewedBy").Preload("Tags").Where("ref = ?", ref).First(&site).Error
	return &site, err
}

func FindDistrictSiteByURL(url string) (*models.DistrictSite, error) {
	var site models.DistrictSite
	err := database.DB.Where("url = ?", url).First(&site).Error
	return &site, err
}

func ListDistrictSites(pagination meta.Pagination, district string, tag string, search string) ([]models.DistrictSite, int64) {
	var sites []models.DistrictSite
	var total int64

	query := database.DB.Model(&models.DistrictSite{}).Where("status = ?", "approved")
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if tag != "" {
		query = query.Where("id IN (SELECT district_site_id FROM district_site_tags INNER JOIN district_tags ON district_tags.id = district_site_tags.district_tag_id WHERE district_tags.name = ?)", strings.ToLower(tag))
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	pagination.Apply(query.Order("created_at desc")).Preload("Submitter").Preload("Tags").Find(&sites)

	return sites, total
}

func ListDistrictSiteRequests(pagination meta.Pagination, status string) ([]models.DistrictSite, int64) {
	var sites []models.DistrictSite
	var total int64

	query := database.DB.Model(&models.DistrictSite{})
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status IN ?", []string{"pending", "hold"})
	}

	query.Count(&total)
	pagination.Apply(query.Order("created_at asc")).Preload("Submitter").Preload("ReviewedBy").Preload("Tags").Find(&sites)

	return sites, total
}

func ListAllDistrictSites(pagination meta.Pagination, district string, search string) ([]models.DistrictSite, int64) {
	var sites []models.DistrictSite
	var total int64

	query := database.DB.Model(&models.DistrictSite{}).Where("status = ?", "approved")
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR url ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	pagination.Apply(query.Order("created_at desc")).Preload("Submitter").Preload("ReviewedBy").Preload("Tags").Find(&sites)

	return sites, total
}

func CountPendingDistrictSites() int64 {
	var count int64
	database.DB.Model(&models.DistrictSite{}).Where("status IN ?", []string{"pending", "hold"}).Count(&count)
	return count
}

func CountApprovedSitesByDistrict(slug string) int64 {
	var count int64
	database.DB.Model(&models.DistrictSite{}).Where("district = ? AND status = ?", slug, "approved").Count(&count)
	return count
}

func FindOrCreateTag(name string) (*models.DistrictTag, error) {
	var tag models.DistrictTag
	normalized := strings.ToLower(strings.TrimSpace(name))
	err := database.DB.Where("name = ?", normalized).FirstOrCreate(&tag, models.DistrictTag{Name: normalized}).Error
	return &tag, err
}

func ReplaceDistrictSiteTags(site *models.DistrictSite, tags []models.DistrictTag) error {
	return database.DB.Model(site).Association("Tags").Replace(tags)
}

func ListPopularTags(limit int) []models.DistrictTag {
	var tags []models.DistrictTag
	database.DB.Raw(`
		SELECT dt.*, COUNT(dst.district_site_id) as count
		FROM district_tags dt
		INNER JOIN district_site_tags dst ON dst.district_tag_id = dt.id
		INNER JOIN district_sites ds ON ds.id = dst.district_site_id AND ds.status = 'approved'
		GROUP BY dt.id
		ORDER BY count DESC
		LIMIT ?
	`, limit).Scan(&tags)
	return tags
}

func ListApprovedSitesWithoutThumbnail() []models.DistrictSite {
	var sites []models.DistrictSite
	database.DB.Where("status = ? AND thumbnail_url = ?", "approved", "").Find(&sites)
	return sites
}

func ListApprovedSitesForThumbnail() []models.DistrictSite {
	var sites []models.DistrictSite
	database.DB.Where("status = ?", "approved").Find(&sites)
	return sites
}