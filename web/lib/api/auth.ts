import { apiClient } from "./client";
import type { AuthUser, LoginResponse } from "@/types/auth";

export async function login(username: string, password: string) {
  return apiClient.post<LoginResponse>("/auth/login", {
    username,
    password,
  });
}

export async function getMe(token: string) {
  return apiClient.get<AuthUser>("/auth/me", token);
}

export type ChangePasswordInput = {
  current_password: string;
  new_password: string;
};

export async function changePassword(token: string, input: ChangePasswordInput) {
  return apiClient.post<{ changed: boolean }>(
    "/auth/change-password",
    input,
    token
  );
}