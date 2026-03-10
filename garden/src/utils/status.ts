import type { AdminUser } from "../types/admin";

export function statusBadge(user: AdminUser | null): string {
  if (!user) return "";
  if (user.account_banned) return "banned";
  if (user.account_disabled) return "disabled";
  if (!user.email_verified) return "unverified";
  return "active";
}