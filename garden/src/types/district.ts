export interface District {
  slug: string;
  name: string;
  description: string;
  background: string;
  foreground: string;
  detail: string;
  site_count: number;
}

export interface CitizenSummary {
  username: string;
  display_name: string;
  avatar_url: string;
}

export interface DistrictSite {
  ref: string;
  district: string;
  district_slug: string;
  title: string;
  url: string;
  description: string;
  thumbnail_url: string;
  tags: string[];
  submitter: CitizenSummary;
  created_at: string;
}

export interface SiteRequest extends DistrictSite {
  status: string;
  reviewed_by: CitizenSummary | null;
  reviewed_at: string | null;
}

export interface AdminSite extends DistrictSite {
  status: string;
  reviewed_by: CitizenSummary | null;
  reviewed_at: string | null;
}