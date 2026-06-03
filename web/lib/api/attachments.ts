import { apiClient } from "./client";
import type {
  AttachmentDownloadURLResponse,
  LetterAttachment,
} from "@/types/attachment";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

export async function listLetterAttachments(token: string, letterID: string) {
  return apiClient.get<LetterAttachment[]>(
    `/letters/${letterID}/attachments`,
    token
  );
}

export async function getAttachmentDownloadURL(
  token: string,
  letterID: string,
  attachmentID: string
) {
  return apiClient.get<AttachmentDownloadURLResponse>(
    `/letters/${letterID}/attachments/${attachmentID}/download-url`,
    token
  );
}

export async function deleteLetterAttachment(
  token: string,
  letterID: string,
  attachmentID: string
) {
  return apiClient.delete<{ deleted: boolean }>(
    `/letters/${letterID}/attachments/${attachmentID}`,
    token
  );
}

export async function uploadLetterAttachments(
  token: string,
  letterID: string,
  files: File[]
) {
  const formData = new FormData();

  files.forEach((file) => {
    formData.append("files", file);
  });

  const res = await fetch(`${API_BASE_URL}/letters/${letterID}/attachments`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: formData,
  });

  const payload = await res.json();

  if (!res.ok || !payload.success) {
    throw new Error(payload.error?.message || "Attachment upload failed");
  }

  return payload.data as LetterAttachment[];
}