import type { Role } from "./auth";

export type User = {
  id: string;
  username: string;
  full_name: string;
  role: Role;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type ListUsersResponse = {
  items: User[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
};