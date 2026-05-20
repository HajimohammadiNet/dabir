import { apiClient } from "./client";
import type { ListAuditLogsResponse } from "@/types/audit";

export type ListAuditLogsParams = {
  page?: number;
  page_size?: number;
  action?: string;
  entity_type?: string;
  actor_user_id?: string;
};

function buildQuery(params: ListAuditLogsParams) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") return;
    searchParams.set(key, String(value));
  });

  const query = searchParams.toString();

  return query ? `?${query}` : "";
}

export async function listAuditLogs(
  token: string,
  params: ListAuditLogsParams = {}
) {
  return apiClient.get<ListAuditLogsResponse>(
    `/audit-logs/${buildQuery(params)}`,
    token
  );
}