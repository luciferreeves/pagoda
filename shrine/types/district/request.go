package district

type SubmitSiteRequest struct {
	District    string   `json:"district"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type ReviewSiteRequest struct {
	Status string `json:"status"`
}

type EditSiteRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	District    *string  `json:"district"`
	Tags        []string `json:"tags"`
}