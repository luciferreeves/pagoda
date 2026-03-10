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

export interface AuditLogEntry {
  system_ref: string;
  actor: string;
  action: string;
  target_type: string;
  target_ref: string;
  summary: string;
  created_at: string;
}

export interface AuditLogDetail extends AuditLogEntry {
  details: string;
}

export const AUDIT_ACTION_LABELS: Record<string, string> = {
  "user.ban": "Ban User",
  "user.unban": "Unban User",
  "user.disable": "Disable User",
  "user.enable": "Enable User",
  "user.role_change": "Role Change",
  "user.edit": "Edit User",
  "user.warn": "Warn User",
  "user.unwarn": "Deactivate Warning",
  "ticket.update": "Update Ticket",
  "district.review": "Review Site",
  "district.edit": "Edit Site",
};

export const AUDIT_TARGET_LABELS: Record<string, string> = {
  user: "User",
  ticket: "Ticket",
  site: "Site",
};