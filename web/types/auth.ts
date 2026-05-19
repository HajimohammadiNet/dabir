export type Role = "superuser" | "editor" | "readonly";

export type AuthUser = {
  id: string;
  username: string;
  full_name: string;
  role: Role;
  is_active?: boolean;
};

export type LoginResponse = {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: AuthUser;
};