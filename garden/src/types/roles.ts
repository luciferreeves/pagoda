export const UserRole = {
  Member: "member",
  Moderator: "moderator",
  Admin: "admin",
} as const;

export type UserRole = (typeof UserRole)[keyof typeof UserRole];