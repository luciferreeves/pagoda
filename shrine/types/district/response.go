package district

import (
	"shrine/types/user"
	"time"
)

type DistrictResponse struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Background  string `json:"background"`
	Foreground  string `json:"foreground"`
	Detail      string `json:"detail"`
	SiteCount   int64  `json:"site_count"`
}

type TagResponse struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type SiteResponse struct {
	Ref          string                      `json:"ref"`
	District     string                      `json:"district"`
	DistrictSlug string                      `json:"district_slug"`
	Title        string                      `json:"title"`
	URL          string                      `json:"url"`
	Description  string                      `json:"description"`
	ThumbnailURL string                      `json:"thumbnail_url"`
	Tags         []string                    `json:"tags"`
	Submitter    user.CitizenSummaryResponse `json:"submitter"`
	CreatedAt    time.Time                   `json:"created_at"`
}

type SiteRequestResponse struct {
	SiteResponse
	Status     string                       `json:"status"`
	ReviewedBy *user.CitizenSummaryResponse `json:"reviewed_by"`
	ReviewedAt *time.Time                   `json:"reviewed_at"`
}

type AdminSiteResponse struct {
	SiteResponse
	Status     string                       `json:"status"`
	ReviewedBy *user.CitizenSummaryResponse `json:"reviewed_by"`
	ReviewedAt *time.Time                   `json:"reviewed_at"`
}