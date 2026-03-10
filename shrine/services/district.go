package services

import (
	"fmt"
	"net/url"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/district"
	"shrine/types/hypertext"
	"shrine/utils/districts"
	"shrine/utils/meta"
	"strings"
	"time"
)

func ListDistricts() []district.DistrictResponse {
	responses := make([]district.DistrictResponse, len(districts.All))
	for index, entry := range districts.All {
		responses[index] = district.DistrictResponse{
			Slug:        string(entry.Slug),
			Name:        entry.Name,
			Description: entry.Description,
			Background:  entry.Background,
			Foreground:  entry.Foreground,
			Detail:      entry.Detail,
			SiteCount:   repositories.CountApprovedSitesByDistrict(string(entry.Slug)),
		}
	}
	return responses
}

func ListDistrictSites(pagination meta.Pagination, slug string, tag string, search string) ([]district.SiteResponse, int64) {
	sites, total := repositories.ListDistrictSites(pagination, slug, tag, search)
	responses := make([]district.SiteResponse, len(sites))
	for index, site := range sites {
		responses[index] = site.ToResponse()
	}
	return responses, total
}

func SubmitSite(userID uint, request district.SubmitSiteRequest) (*district.SiteRequestResponse, *hypertext.ServiceError) {
	if !districts.IsValid(request.District) {
		return nil, fail(enums.BadRequest, messages.InvalidDistrict)
	}

	title := strings.TrimSpace(request.Title)
	if title == "" || len(title) > 200 {
		return nil, fail(enums.BadRequest, messages.SiteTitleRequired)
	}

	siteURL := strings.TrimSpace(request.URL)
	if !isValidURL(siteURL) {
		return nil, fail(enums.BadRequest, messages.SiteURLRequired)
	}

	description := strings.TrimSpace(request.Description)
	if len(description) > 1000 {
		return nil, fail(enums.BadRequest, messages.SiteDescriptionLong)
	}

	tags, serviceErr := validateTags(request.Tags)
	if serviceErr != nil {
		return nil, serviceErr
	}

	existing, _ := repositories.FindDistrictSiteByURL(siteURL)
	if existing != nil {
		return nil, fail(enums.Conflict, messages.SiteURLTaken)
	}

	site := models.DistrictSite{
		District:    request.District,
		SubmitterID: userID,
		Title:       title,
		URL:         siteURL,
		Description: description,
		Status:      string(enums.SitePending),
	}

	if err := repositories.CreateDistrictSite(&site); err != nil {
		return nil, fail(enums.Internal, messages.FailedSubmitSite)
	}

	if len(tags) > 0 {
		repositories.ReplaceDistrictSiteTags(&site, tags)
	}

	site.Submitter = models.User{}
	record, _ := repositories.FindDistrictSiteByRef(site.Ref)
	response := record.ToRequestResponse()
	return &response, nil
}

func ReviewSite(adminID uint, ref string, request district.ReviewSiteRequest) (*district.SiteRequestResponse, *hypertext.ServiceError) {
	site, serviceErr := resolveDistrictSite(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	status := enums.SiteStatus(request.Status)
	switch status {
	case enums.SiteApproved, enums.SiteDenied, enums.SiteHold:
	default:
		return nil, fail(enums.BadRequest, messages.InvalidSiteStatus)
	}

	site.Status = request.Status
	site.ReviewedByID = &adminID
	now := time.Now()
	site.ReviewedAt = &now

	if err := repositories.UpdateDistrictSite(site); err != nil {
		return nil, fail(enums.Internal, messages.FailedReviewSite)
	}

	var auditMessage string
	switch status {
	case enums.SiteApproved:
		auditMessage = fmt.Sprintf(messages.AuditApprovedSite, site.Ref)
	case enums.SiteDenied:
		auditMessage = fmt.Sprintf(messages.AuditDeniedSite, site.Ref)
	case enums.SiteHold:
		auditMessage = fmt.Sprintf(messages.AuditHeldSite, site.Ref)
	}

	repositories.LogAction(adminID, "district.review", "site", site.Ref, auditMessage, request)

	if status == enums.SiteApproved {
		go GenerateSiteThumbnail(site.Ref, site.URL)
	}

	record, _ := repositories.FindDistrictSiteByRef(site.Ref)
	response := record.ToRequestResponse()
	return &response, nil
}

func EditSite(adminID uint, ref string, request district.EditSiteRequest) (*district.AdminSiteResponse, *hypertext.ServiceError) {
	site, serviceErr := resolveDistrictSite(ref)
	if serviceErr != nil {
		return nil, serviceErr
	}

	if request.Title != nil {
		title := strings.TrimSpace(*request.Title)
		if title == "" || len(title) > 200 {
			return nil, fail(enums.BadRequest, messages.SiteTitleRequired)
		}
		site.Title = title
	}

	if request.Description != nil {
		description := strings.TrimSpace(*request.Description)
		if len(description) > 1000 {
			return nil, fail(enums.BadRequest, messages.SiteDescriptionLong)
		}
		site.Description = description
	}

	if request.District != nil {
		if !districts.IsValid(*request.District) {
			return nil, fail(enums.BadRequest, messages.InvalidDistrict)
		}
		site.District = *request.District
	}

	if request.Tags != nil {
		tags, tagErr := validateTags(request.Tags)
		if tagErr != nil {
			return nil, tagErr
		}
		repositories.ReplaceDistrictSiteTags(site, tags)
	}

	if err := repositories.UpdateDistrictSite(site); err != nil {
		return nil, fail(enums.Internal, messages.FailedEditSite)
	}

	repositories.LogAction(adminID, "district.edit", "site", site.Ref, fmt.Sprintf(messages.AuditEditedSite, site.Ref), request)

	record, _ := repositories.FindDistrictSiteByRef(site.Ref)
	response := record.ToAdminResponse()
	return &response, nil
}

func ListSiteRequests(pagination meta.Pagination, status string) ([]district.SiteRequestResponse, int64) {
	sites, total := repositories.ListDistrictSiteRequests(pagination, status)
	responses := make([]district.SiteRequestResponse, len(sites))
	for index, site := range sites {
		responses[index] = site.ToRequestResponse()
	}
	return responses, total
}

func ListAdminSites(pagination meta.Pagination, slug string, search string) ([]district.AdminSiteResponse, int64) {
	sites, total := repositories.ListAllDistrictSites(pagination, slug, search)
	responses := make([]district.AdminSiteResponse, len(sites))
	for index, site := range sites {
		responses[index] = site.ToAdminResponse()
	}
	return responses, total
}

func CountPendingSites() int64 {
	return repositories.CountPendingDistrictSites()
}

func resolveDistrictSite(ref string) (*models.DistrictSite, *hypertext.ServiceError) {
	site, err := repositories.FindDistrictSiteByRef(ref)
	if err != nil {
		return nil, fail(enums.NotFound, messages.SiteNotFound)
	}
	return site, nil
}

func isValidURL(raw string) bool {
	if raw == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func validateTags(raw []string) ([]models.DistrictTag, *hypertext.ServiceError) {
	if len(raw) > 5 {
		return nil, fail(enums.BadRequest, messages.TooManyTags)
	}

	var tags []models.DistrictTag
	for _, name := range raw {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		if len(trimmed) > 50 {
			return nil, fail(enums.BadRequest, messages.TagTooLong)
		}
		tag, err := repositories.FindOrCreateTag(trimmed)
		if err != nil {
			return nil, fail(enums.Internal, messages.FailedSubmitSite)
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}