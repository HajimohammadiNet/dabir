import { apiClient } from "./client";
import type { PublicSettings } from "@/types/settings";

export async function getPublicSettings() {
  return apiClient.get<PublicSettings>("/settings/public");
}