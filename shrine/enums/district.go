package enums

type DistrictSlug string

const (
	Arcadia           DistrictSlug = "arcadia"
	Arles             DistrictSlug = "arles"
	Hollywood         DistrictSlug = "hollywood"
	Oxford            DistrictSlug = "oxford"
	Petsburg          DistrictSlug = "petsburg"
	Purgatory         DistrictSlug = "purgatory"
	SiliconValley     DistrictSlug = "silicon-valley"
	SilverLake        DistrictSlug = "silver-lake"
	StratfordUponAvon DistrictSlug = "stratford-upon-avon"
	Tokyo             DistrictSlug = "tokyo"
)

type SiteStatus string

const (
	SitePending  SiteStatus = "pending"
	SiteApproved SiteStatus = "approved"
	SiteDenied   SiteStatus = "denied"
	SiteHold     SiteStatus = "hold"
)