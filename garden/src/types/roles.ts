export const UserRole = {
  Member: "member",
  Moderator: "moderator",
  Admin: "admin",
  Owner: "owner",
} as const;

export type UserRole = (typeof UserRole)[keyof typeof UserRole];

export const ROLE_LABELS: Record<string, string> = {
  member: "Member",
  moderator: "Moderator",
  admin: "Admin",
  owner: "Owner",
};