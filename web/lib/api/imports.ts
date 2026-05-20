import { apiClient } from "./client";
import type { CommitImportResponse, ImportJob } from "@/types/import";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

export async function previewLettersImport(token: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(`${API_BASE_URL}/imports/letters/preview`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: formData,
  });

  const payload = await res.json();

  if (!res.ok || !payload.success) {
    throw new Error(payload.error?.message || "Import preview failed");
  }

  return payload.data as ImportJob;
}

export async function getImportJob(token: string, id: string) {
  return apiClient.get<ImportJob>(`/imports/${id}`, token);
}

export async function commitLettersImport(token: string, id: string) {
  return apiClient.post<CommitImportResponse>(
    `/imports/letters/${id}/commit`,
    undefined,
    token
  );
}