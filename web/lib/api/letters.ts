import { apiClient } from "./client";
import type {
  Letter,
  LetterDirection,
  ListLettersResponse,
} from "@/types/letter";

export type ListLettersParams = {
  page?: number;
  page_size?: number;
  search?: string;
  direction?: LetterDirection;
  sender?: string;
  receiver?: string;
  registrar_name?: string;
  from_date?: string;
  to_date?: string;
  include_deleted?: boolean;
  sort_by?: "created_at" | "letter_date" | "letter_number";
  sort_order?: "asc" | "desc";
};

function buildQuery(params: Record<string, unknown>) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    searchParams.set(key, String(value));
  });

  const query = searchParams.toString();

  return query ? `?${query}` : "";
}

export async function listLetters(
  token: string,
  params: ListLettersParams = {}
) {
  return apiClient.get<ListLettersResponse>(
    `/letters/${buildQuery(params)}`,
    token
  );
}

export type CreateLetterInput = {
  direction?: LetterDirection;
  display_letter_number?: string | null;
  title: string;
  letter_date: string;
  sender: string;
  receiver: string;
  description?: string | null;
};

export async function createLetter(token: string, input: CreateLetterInput) {
  return apiClient.post<Letter>("/letters/", input, token);
}

export async function getLetter(token: string, id: string) {
  return apiClient.get<Letter>(`/letters/${id}`, token);
}

export type UpdateLetterInput = {
  display_letter_number?: string | null;
  title: string;
  letter_date: string;
  sender: string;
  receiver: string;
  description?: string | null;
};

export async function updateLetter(
  token: string,
  id: string,
  input: UpdateLetterInput
) {
  return apiClient.patch<Letter>(`/letters/${id}`, input, token);
}

export async function deleteLetter(token: string, id: string) {
  return apiClient.delete<{ deleted: boolean }>(`/letters/${id}`, token);
}

export type LetterNumberSuggestion = {
  prefix: string;
  last_number?: string | null;
  suggested_number?: string | null;
};

export async function getLetterNumberSuggestion(
  token: string,
  prefix: string,
  direction: LetterDirection = "incoming"
) {
  const query = buildQuery({ prefix, direction });

  return apiClient.get<LetterNumberSuggestion>(
    `/letters/number-suggestion${query}`,
    token
  );
}
