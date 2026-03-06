export const UserRole = {
  Member: "member",
  Moderator: "moderator",
  Admin: "admin",
  Owner: "owner",
} as const;

export type UserRole = (typeof UserRole)[keyof typeof UserRole];