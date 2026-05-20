import { apiClient } from "./client";
import type {
  InitializeSetupInput,
  InitializeSetupResponse,
  SetupStatus,
} from "@/types/setup";

export async function getSetupStatus() {
  return apiClient.get<SetupStatus>("/setup/status");
}

export async function initializeSetup(input: InitializeSetupInput) {
  return apiClient.post<InitializeSetupResponse>("/setup/initialize", input);
}