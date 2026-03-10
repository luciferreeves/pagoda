export interface AdminUser {
  username: string;
  email: string;
  display_name: string;
  bio: string;
  birthday: string | null;
  avatar_url: string;
  blinkie_url: string;
  website: string;
  location: string;
  pronouns: string;
  signature: string;
  role: string;
  jade: number;
  honor: number;
  email_verified: boolean;
  warning_count: number;
  account_banned: boolean;
  banned_reason: string;
  banned_at: string | null;
  account_disabled: boolean;
  disabled_reason: string;
  disabled_at: string | null;
  disabled_until: string | null;
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