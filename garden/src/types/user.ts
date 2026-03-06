import type { UserRole } from "./roles";

export interface User {
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
  role: UserRole;
  created_at: string;
}