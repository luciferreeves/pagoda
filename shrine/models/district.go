package models

import (
	"shrine/enums"
	"shrine/types/district"
	"shrine/utils/crypto"
	"shrine/utils/districts"
	"shrine/utils/storage"
	"time"

	"gorm.io/gorm"
)

type DistrictSite struct {
	gorm.Model
	Ref          string        `gorm:"size:12;uniqueIndex;not null"`
	District     string        `gorm:"size:30;index;not null"`
	SubmitterID  uint          `gorm:"index;not null"`
	Submitter    User          `gorm:"foreignKey:SubmitterID"`
	Title        string        `gorm:"size:200;not null"`
	URL          string        `gorm:"size:500;uniqueIndex;not null"`
	Description  string        `gorm:"size:1000"`
	ThumbnailURL string        `gorm:"size:500"`
	Status       string        `gorm:"size:10;not null;default:pending;index"`
	ReviewedByID *uint         `gorm:"index"`
	ReviewedBy   *User         `gorm:"foreignKey:ReviewedByID"`
	ReviewedAt   *time.Time
	Tags         []DistrictTag `gorm:"many2many:district_site_tags"`
}

type DistrictTag struct {
	gorm.Model
	Name string `gorm:"size:50;uniqueIndex;not null"`
}

func (self *DistrictSite) BeforeCreate(tx *gorm.DB) error {
	if self.Ref == "" {
		self.Ref = crypto.Ref()
	}
	return nil
}

func (self *DistrictSite) ToResponse() district.SiteResponse {
	tagNames := make([]string, len(self.Tags))
	for index, tag := range self.Tags {
		tagNames[index] = tag.Name
	}

	districtInfo, _ := districts.Find(enums.DistrictSlug(self.District))

	return district.SiteResponse{
		Ref:          self.Ref,
		District:     districtInfo.Name,
		DistrictSlug: self.District,
		Title:        self.Title,
		URL:          self.URL,
		Description:  self.Description,
		ThumbnailURL: storage.ResolveCDN(self.ThumbnailURL),
		Tags:         tagNames,
		Submitter:    self.Submitter.ToSummary(),
		CreatedAt:    self.CreatedAt,
	}
}

func (self *DistrictSite) ToRequestResponse() district.SiteRequestResponse {
	response := district.SiteRequestResponse{
		SiteResponse: self.ToResponse(),
		Status:       self.Status,
	}

	if self.ReviewedBy != nil {
		summary := self.ReviewedBy.ToSummary()
		response.ReviewedBy = &summary
	}

	if self.ReviewedAt != nil {
		response.ReviewedAt = self.ReviewedAt
	}

	return response
}

func (self *DistrictSite) ToAdminResponse() district.AdminSiteResponse {
	response := district.AdminSiteResponse{
		SiteResponse: self.ToResponse(),
		Status:       self.Status,
	}

	if self.ReviewedBy != nil {
		summary := self.ReviewedBy.ToSummary()
		response.ReviewedBy = &summary
	}

	if self.ReviewedAt != nil {
		response.ReviewedAt = self.ReviewedAt
	}

	return response
}