export interface CitizenSummary {
  username: string;
  display_name: string;
  avatar_url: string;
}

export interface Stats {
  citizens: number;
  online: number;
  unread_letters: number;
  pending_districts: number;
  newest_citizens: CitizenSummary[];
  online_citizens: CitizenSummary[];
}