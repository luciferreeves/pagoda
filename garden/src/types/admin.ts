export interface AdminUser {
  username: string;
  email: string;
  display_name: string;
  avatar_url: string;
  role: string;
  email_verified: boolean;
  account_banned: boolean;
  banned_reason: string;
  banned_at: string | null;
  account_disabled: boolean;
  disabled_reason: string;
  disabled_at: string | null;
  last_seen_at: string | null;
  created_at: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}