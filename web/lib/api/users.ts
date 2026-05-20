import { apiClient } from "./client";
import type { Role } from "@/types/auth";
import type { ListUsersResponse, User } from "@/types/user";

export type ListUsersParams = {
  page?: number;
  page_size?: number;
  search?: string;
  role?: Role;
  is_active?: boolean;
};

function buildQuery(params: ListUsersParams) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    searchParams.set(key, String(value));
  });

  const query = searchParams.toString();

  return query ? `?${query}` : "";
}

export async function listUsers(token: string, params: ListUsersParams = {}) {
  return apiClient.get<ListUsersResponse>(`/users/${buildQuery(params)}`, token);
}

export type CreateUserInput = {
  username: string;
  full_name: string;
  password: string;
  role: Role;
};

export async function createUser(token: string, input: CreateUserInput) {
  return apiClient.post<User>("/users/", input, token);
}

export type UpdateUserInput = {
  full_name: string;
  role: Role;
  is_active?: boolean;
};

export async function updateUser(
  token: string,
  id: string,
  input: UpdateUserInput
) {
  return apiClient.patch<User>(`/users/${id}`, input, token);
}

export async function activateUser(token: string, id: string) {
  return apiClient.patch<{ activated: boolean }>(
    `/users/${id}/activate`,
    undefined,
    token
  );
}

export async function deactivateUser(token: string, id: string) {
  return apiClient.patch<{ deactivated: boolean }>(
    `/users/${id}/deactivate`,
    undefined,
    token
  );
}