package messages

const (
	InvalidDistrict       = "Invalid district."
	InvalidSiteStatus     = "Invalid site status."
	SiteNotFound          = "Site not found."
	SiteTitleRequired     = "Title must be between 1 and 200 characters."
	SiteURLRequired       = "A valid URL is required."
	SiteURLTaken          = "This URL has already been submitted."
	SiteDescriptionLong   = "Description must be 1000 characters or fewer."
	TooManyTags           = "A site can have up to 5 tags."
	TagTooLong            = "Each tag must be 50 characters or fewer."
	FailedSubmitSite      = "Failed to submit site."
	FailedReviewSite      = "Failed to review site."
	FailedEditSite        = "Failed to update site."
	SiteAlreadyReviewed   = "This site has already been reviewed."
	AuditApprovedSite     = "Approved site %s"
	AuditDeniedSite       = "Denied site %s"
	AuditHeldSite         = "Put site %s on hold"
	AuditEditedSite       = "Edited site %s"
)